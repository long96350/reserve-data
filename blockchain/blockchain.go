package blockchain

import (
	"fmt"
	"math/big"
	"path/filepath"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/blockchain"
	huobiblockchain "github.com/KyberNetwork/reserve-data/exchange/huobi/blockchain"
	commonv3 "github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/KyberNetwork/reserve-data/reservesetting/storage"
)

const (
	pricingOP = "pricingOP"
	depositOP = "depositOP"
)

var (
	Big0   = big.NewInt(0)
	BigMax = big.NewInt(10).Exp(big.NewInt(10), big.NewInt(33), nil)
)

// tbindex is where the token data stored in blockchain.
// In blockchain, data of a token (sell/buy rates) is stored in an array of 32 bytes values called (tokenRatesCompactData).
// Each data is stored in a byte.
// https://github.com/KyberNetwork/smart-contracts/blob/fed8e09dc6e4365e1597474d9b3f53634eb405d2/contracts/ConversionRates.sol#L48
type tbindex struct {
	// BulkIndex is the index of bytes32 value that store data of multiple tokens.
	BulkIndex uint64
	// IndexInBulk is the index in the above BulkIndex value where the sell/buy rates are stored following structure:
	// sell: IndexInBulk + 4
	// buy: IndexInBulk + 8
	IndexInBulk uint64
}

// newTBIndex creates new tbindex instance with given parameters.
func newTBIndex(bulkIndex, indexInBulk uint64) tbindex {
	return tbindex{BulkIndex: bulkIndex, IndexInBulk: indexInBulk}
}

type Blockchain struct {
	*blockchain.BaseBlockchain
	wrapper      *blockchain.Contract
	pricing      *blockchain.Contract
	reserve      *blockchain.Contract
	tokenIndices map[string]tbindex
	mu           sync.RWMutex

	localSetRateNonce     uint64
	setRateNonceTimestamp uint64
	contractAddress       *common.ContractAddressConfiguration
	sr                    storage.SettingReader
	l                     *zap.SugaredLogger
}

func (bc *Blockchain) StandardGasPrice() float64 {
	// we use node's recommended gas price because gas station is not returning
	// correct gas price now
	price, err := bc.RecommendedGasPriceFromNode()
	if err != nil {
		return 0
	}
	return common.BigToFloat(price, 9)
}

func (bc *Blockchain) CheckTokenIndices(tokenAddr ethereum.Address) error {
	opts := bc.GetCallOpts(0)
	pricingAddr := bc.contractAddress.Pricing
	tokenAddrs := []ethereum.Address{}
	tokenAddrs = append(tokenAddrs, tokenAddr)
	_, _, err := bc.GeneratedGetTokenIndicies(
		opts,
		pricingAddr,
		tokenAddrs,
	)
	if err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) LoadAndSetTokenIndices(tokenAddrs []ethereum.Address) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.tokenIndices = map[string]tbindex{}
	opts := bc.GetCallOpts(0)
	pricingAddr := bc.contractAddress.Pricing
	bulkIndices, indicesInBulk, err := bc.GeneratedGetTokenIndicies(
		opts,
		pricingAddr,
		tokenAddrs,
	)
	if err != nil {
		return err
	}
	for i, tok := range tokenAddrs {
		bc.tokenIndices[tok.Hex()] = newTBIndex(
			bulkIndices[i].Uint64(),
			indicesInBulk[i].Uint64(),
		)
	}
	bc.l.Infof("Token indices: %+v", bc.tokenIndices)
	return nil
}

func (bc *Blockchain) RegisterPricingOperator(signer blockchain.Signer, nonceCorpus blockchain.NonceCorpus) {
	bc.l.Infow("reserve pricing address", "address", signer.GetAddress().Hex())
	bc.MustRegisterOperator(pricingOP, blockchain.NewOperator(signer, nonceCorpus))
}

func (bc *Blockchain) RegisterDepositOperator(signer blockchain.Signer, nonceCorpus blockchain.NonceCorpus) {
	bc.l.Infow("reserve depositor address", "address", signer.GetAddress().Hex())
	bc.MustRegisterOperator(depositOP, blockchain.NewOperator(signer, nonceCorpus))
}

func readablePrint(data map[ethereum.Address]byte) string {
	result := ""
	for addr, b := range data {
		result = result + "|" + fmt.Sprintf("%s-%d", addr.Hex(), b)
	}
	return result
}

//====================== Write calls ===============================

// TODO: Need better test coverage
// we got a bug when compact is not set to old compact
// or when one of buy/sell got overflowed, it discards
// the other's compact
func (bc *Blockchain) SetRates(
	tokens []ethereum.Address,
	buys []*big.Int,
	sells []*big.Int,
	block *big.Int,
	nonce *big.Int,
	gasPrice *big.Int) (*types.Transaction, error) {
	pricingAddr := bc.contractAddress.Pricing
	block.Add(block, big.NewInt(1))
	copts := bc.GetCallOpts(0)
	baseBuys, baseSells, _, _, _, err := bc.GeneratedGetTokenRates(
		copts, pricingAddr, tokens,
	)
	if err != nil {
		return nil, err
	}

	// This is commented out because we dont want to make too much of change. Don't remove
	// this check, it can be useful in the future.
	//
	// Don't submit any txs if it is just trying to set all tokens to 0 when they are already 0
	// if common.AllZero(buys, sells, baseBuys, baseSells) {
	// 	return nil, errors.New("Trying to set all rates to 0 but they are already 0. Skip the tx.")
	// }

	baseTokens := []ethereum.Address{}
	newBSells := []*big.Int{}
	newBBuys := []*big.Int{}
	newCSells := map[ethereum.Address]byte{}
	newCBuys := map[ethereum.Address]byte{}
	for i, token := range tokens {
		compactSell, overflow1 := BigIntToCompactRate(sells[i], baseSells[i])
		compactBuy, overflow2 := BigIntToCompactRate(buys[i], baseBuys[i])
		if overflow1 || overflow2 {
			baseTokens = append(baseTokens, token)
			newBSells = append(newBSells, sells[i])
			newBBuys = append(newBBuys, buys[i])
			newCSells[token] = 0
			newCBuys[token] = 0
		} else {
			newCSells[token] = compactSell.Compact
			newCBuys[token] = compactBuy.Compact
		}
	}
	bbuys, bsells, indices := BuildCompactBulk(
		newCBuys,
		newCSells,
		bc.tokenIndices,
	)
	opts, err := bc.GetTxOpts(pricingOP, nonce, gasPrice, nil)
	if err != nil {
		bc.l.Infow("Getting transaction opts failed", "err", err)
		return nil, err
	}
	var tx *types.Transaction
	if len(baseTokens) > 0 {
		// set base tx
		tx, err = bc.GeneratedSetBaseRate(
			opts, baseTokens, newBBuys, newBSells,
			bbuys, bsells, block, indices)
		if tx != nil {
			bc.l.Infof(
				"broadcasting setbase tx %s, target buys(%s), target sells(%s), old base buy(%s) || old base sell(%s) || new base buy(%s) || new base sell(%s) || new compact buy(%s) || new compact sell(%s) || new buy bulk(%v) || new sell bulk(%v) || indices(%v)",
				tx.Hash().Hex(),
				buys, sells,
				baseBuys, baseSells,
				newBBuys, newBSells,
				readablePrint(newCBuys), readablePrint(newCSells),
				bbuys, bsells, indices,
			)
		}
	} else {
		// update compact tx
		tx, err = bc.GeneratedSetCompactData(
			opts, bbuys, bsells, block, indices)
		if tx != nil {
			bc.l.Infof(
				"broadcasting setcompact tx %s, target buys(%s), target sells(%s), old base buy(%s) || old base sell(%s) || new compact buy(%s) || new compact sell(%s) || new buy bulk(%v) || new sell bulk(%v) || indices(%v)",
				tx.Hash().Hex(),
				buys, sells,
				baseBuys, baseSells,
				readablePrint(newCBuys), readablePrint(newCSells),
				bbuys, bsells, indices,
			)
		}
	}
	if err != nil {
		return nil, err
	}
	return bc.SignAndBroadcast(tx, pricingOP)

}

func (bc *Blockchain) Send(
	asset commonv3.Asset,
	amount *big.Int,
	dest ethereum.Address) (*types.Transaction, error) {

	opts, err := bc.GetTxOpts(depositOP, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	tx, err := bc.GeneratedWithdraw(
		opts,
		asset.Address,
		amount,
		dest)
	if err != nil {
		return nil, err
	}
	return bc.SignAndBroadcast(tx, depositOP)
}

func (bc *Blockchain) SetImbalanceStepFunction(token ethereum.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	opts, err := bc.GetTxOpts(pricingOP, nil, nil, nil)
	if err != nil {
		bc.l.Infow("Getting transaction opts failed", "err", err)
		return nil, err
	}
	tx, err := bc.GeneratedSetImbalanceStepFunction(opts, token, xBuy, yBuy, xSell, ySell)
	if err != nil {
		return nil, err
	}
	return bc.SignAndBroadcast(tx, pricingOP)
}

func (bc *Blockchain) SetQtyStepFunction(token ethereum.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	opts, err := bc.GetTxOpts(pricingOP, nil, nil, nil)
	if err != nil {
		bc.l.Infow("Getting transaction opts failed", "err", err)
		return nil, err
	}
	tx, err := bc.GeneratedSetQtyStepFunction(opts, token, xBuy, yBuy, xSell, ySell)
	if err != nil {
		return nil, err
	}
	return bc.SignAndBroadcast(tx, pricingOP)
}

//====================== Readonly calls ============================
func (bc *Blockchain) FetchBalanceData(reserve ethereum.Address, atBlock uint64) (map[string]common.BalanceEntry, error) {
	result := map[string]common.BalanceEntry{}
	tokens := []ethereum.Address{}
	assets, err := bc.sr.GetTransferableAssets()
	if err != nil {
		return result, err
	}
	for _, tok := range assets {
		tokens = append(tokens, tok.Address)
	}
	timestamp := common.GetTimestamp()
	opts := bc.GetCallOpts(atBlock)
	balances, err := bc.GeneratedGetBalances(opts, reserve, tokens)
	returnTime := common.GetTimestamp()
	bc.l.Infow("Fetcher ------> balances", "balances", balances, "err", err)
	if err != nil {
		for _, token := range assets {
			// TODO: should store token id instead of symbol
			result[token.Symbol] = common.BalanceEntry{
				Valid:      false,
				Error:      err.Error(),
				Timestamp:  timestamp,
				ReturnTime: returnTime,
			}
		}
	} else {
		for i, tok := range assets {
			if balances[i].Cmp(Big0) == 0 || balances[i].Cmp(BigMax) > 0 {
				// TODO: log asset_id instead of token
				bc.l.Infow("Fetcher ------> balances of token is invalid", "token_symbol", tok.Symbol)
				// TODO: should store token id instead of symbol
				result[tok.Symbol] = common.BalanceEntry{
					Valid:      false,
					Error:      "Got strange balances from node. It equals to 0 or is bigger than 10^33",
					Timestamp:  timestamp,
					ReturnTime: returnTime,
					Balance:    common.RawBalance(*balances[i]),
				}
			} else {
				// TODO: should store token id instead of symbol
				result[tok.Symbol] = common.BalanceEntry{
					Valid:      true,
					Timestamp:  timestamp,
					ReturnTime: returnTime,
					Balance:    common.RawBalance(*balances[i]),
				}
			}
		}
	}
	return result, nil
}

func (bc *Blockchain) FetchRates(atBlock uint64, currentBlock uint64) (common.AllRateEntry, error) {
	result := common.AllRateEntry{}
	tokenAddrs := []ethereum.Address{}
	validTokens := []commonv3.Asset{}
	assets, err := bc.sr.GetTransferableAssets()
	if err != nil {
		return result, err
	}
	for _, s := range assets {
		// TODO: add a isETH method
		if s.Symbol != "ETH" {
			tokenAddrs = append(tokenAddrs, s.Address)
			validTokens = append(validTokens, s)
		}
	}
	timestamp := common.GetTimestamp()
	opts := bc.GetCallOpts(atBlock)
	pricingAddr := bc.contractAddress.Pricing
	baseBuys, baseSells, compactBuys, compactSells, blocks, err := bc.GeneratedGetTokenRates(
		opts, pricingAddr, tokenAddrs,
	)
	if err != nil {
		return result, err
	}
	returnTime := common.GetTimestamp()
	result.Timestamp = timestamp
	result.ReturnTime = returnTime
	result.BlockNumber = currentBlock

	result.Data = map[uint64]common.RateEntry{}
	for i, token := range validTokens {
		result.Data[token.ID] = common.NewRateEntry(
			baseBuys[i],
			compactBuys[i],
			baseSells[i],
			compactSells[i],
			blocks[i].Uint64(),
		)
	}
	return result, nil
}

func (bc *Blockchain) GetPrice(token ethereum.Address, block *big.Int, priceType string, qty *big.Int, atBlock uint64) (*big.Int, error) {
	opts := bc.GetCallOpts(atBlock)
	if priceType == "buy" {
		return bc.GeneratedGetRate(opts, token, block, true, qty)
	}
	return bc.GeneratedGetRate(opts, token, block, false, qty)
}

// SetRateMinedNonce returns nonce of the pricing operator in confirmed
// state (not pending state).
//
// Getting mined nonce is not simple because there might be lag between
// node leading us to get outdated mined nonce from an unsynced node.
// To overcome this situation, we will keep a local nonce and require
// the nonce from node to be equal or greater than it.
// If the nonce from node is smaller than the local one, we will use
// the local one. However, if the local one stay with the same value
// for more than 15 mins, the local one is considered incorrect
// because the chain might be reorganized so we will invalidate it
// and assign it to the nonce from node.
func (bc *Blockchain) SetRateMinedNonce() (uint64, error) {
	const localNonceExpiration = time.Minute * 2
	nonceFromNode, err := bc.GetMinedNonce(pricingOP)
	if err != nil {
		return nonceFromNode, err
	}
	if nonceFromNode < bc.localSetRateNonce {
		bc.l.Infof("SET_RATE_MINED_NONCE: nonce returned from node %d is smaller than cached nonce: %d",
			nonceFromNode, bc.localSetRateNonce)
		if common.NowInMillis()-bc.setRateNonceTimestamp > uint64(localNonceExpiration/time.Millisecond) {
			bc.l.Infof("SET_RATE_MINED_NONCE: cached nonce %d stalled, overwriting with nonce from node %d",
				bc.localSetRateNonce, nonceFromNode)
			bc.localSetRateNonce = nonceFromNode
			bc.setRateNonceTimestamp = common.NowInMillis()
			return nonceFromNode, nil
		}
		bc.l.Infof("SET_RATE_MINED_NONCE: using cached nonce %d instead of nonce from node %d",
			bc.localSetRateNonce, nonceFromNode)
		return bc.localSetRateNonce, nil
	}

	bc.l.Infof("SET_RATE_MINED_NONCE: updating cached nonce, current: %d, new: %d",
		bc.localSetRateNonce, nonceFromNode)
	bc.localSetRateNonce = nonceFromNode
	bc.setRateNonceTimestamp = common.NowInMillis()
	return nonceFromNode, nil
}

func NewBlockchain(base *blockchain.BaseBlockchain,
	contractAddressConf *common.ContractAddressConfiguration,
	sr storage.SettingReader,
	l *zap.SugaredLogger,
) (*Blockchain, error) {
	l.Infow("wrapper address", "address", contractAddressConf.Wrapper.Hex())
	wrapper := blockchain.NewContract(
		contractAddressConf.Wrapper,
		filepath.Join(common.CurrentDir(), "wrapper.abi"),
	)

	l.Infow("reserve address", "address", contractAddressConf.Reserve.Hex())
	reserve := blockchain.NewContract(
		contractAddressConf.Reserve,
		filepath.Join(common.CurrentDir(), "reserve.abi"),
	)

	l.Infow("pricing address", "address", contractAddressConf.Pricing.Hex())
	pricing := blockchain.NewContract(
		contractAddressConf.Pricing,
		filepath.Join(common.CurrentDir(), "pricing.abi"),
	)

	return &Blockchain{
		BaseBlockchain:  base,
		wrapper:         wrapper,
		pricing:         pricing,
		reserve:         reserve,
		contractAddress: contractAddressConf,
		sr:              sr,
		l:               l,
	}, nil
}

// GetPricingOPAddress return pricing op address
func (bc *Blockchain) GetPricingOPAddress() ethereum.Address {
	return bc.MustGetOperator(pricingOP).Address
}

// GetDepositOPAddress return deposit op address
func (bc *Blockchain) GetDepositOPAddress() ethereum.Address {
	return bc.MustGetOperator(depositOP).Address
}

// GetIntermediatorOPAddress return intermediator op address
func (bc *Blockchain) GetIntermediatorOPAddress() ethereum.Address {
	return bc.MustGetOperator(huobiblockchain.HuobiOP).Address
}

// GetWrapperAddress return wrapper address
func (bc *Blockchain) GetWrapperAddress() ethereum.Address {
	return bc.contractAddress.Wrapper
}

// GetInternalNetworkAddress return internal network address
func (bc *Blockchain) GetInternalNetworkAddress() ethereum.Address {
	return bc.contractAddress.InternalNetwork
}

// GetNetworkAddress return network address
func (bc *Blockchain) GetNetworkAddress() ethereum.Address {
	return bc.contractAddress.Network
}
