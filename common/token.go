package common

// TokenExchangeSetting contains necessary information on exchange to List a token on the fly
type TokenExchangeSetting struct {
	DepositAddress string       `json:"deposit_address"`
	Info           ExchangeInfo `json:"exchange_info"`
	Fee            TokenFee     `json:"fee"`
	MinDeposit     float64      `json:"min_deposit"`
}

type TokenUpdate struct {
	Token       Token                           `json:"token"`
	Exchanges   map[string]TokenExchangeSetting `json:"exchanges"`
	PWIEq       PWIEquationTokenV2              `json:"pwis_equation"`
	TargetQty   TargetQtyV2                     `json:"target_qty"`
	QuadraticEq RebalanceQuadraticEquation      `json:"rebalance_quadratic"`
}

type TokenFee struct {
	Withdraw float64 `json:"withdraw"`
	Deposit  float64 `json:"deposit"`
}

type Token struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Address              string `json:"address"`
	Decimals             int64  `json:"decimals"`
	Internal             bool   `json:"internal"`
	LastActivationChange uint64 `json:"last_activation_change"`
	CreationTime         uint64 `json:"creation_time"`
}

// NewToken creates a new Token.
func NewToken(id, name, address string, decimal int64, active, internal bool, timepoint uint64) Token {
	return Token{
		ID:                   id,
		Name:                 name,
		Address:              address,
		Decimals:             decimal,
		Internal:             internal,
		LastActivationChange: timepoint,
	}
}

func (t Token) IsETH() bool {
	return t.ID == "ETH"
}
