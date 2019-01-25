package exchange

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// BittrexInterface contains the methods to interact with Bittrex centralized exchange.
type BittrexInterface interface {
	FetchOnePairData(baseID, quoteID string) (Bittresp, error)

	GetInfo() (Bittinfo, error)

	GetExchangeInfo() (BittExchangeInfo, error)

	GetDepositAddress(currency string) (BittrexDepositAddress, error)

	GetAccountTradeHistory(base, quote common.Token) (BittTradeHistory, error)

	Withdraw(
		token common.Token,
		amount *big.Int,
		address ethereum.Address) (Bittwithdraw, error)

	Trade(
		tradeType string,
		base, quote common.Token,
		rate, amount float64) (Bitttrade, error)

	CancelOrder(uuid string) (Bittcancelorder, error)

	DepositHistory(currency string) (Bittdeposithistory, error)

	WithdrawHistory(currency string) (Bittwithdrawhistory, error)

	OrderStatus(uuid string) (Bitttraderesult, error)
}
