package common

type WsFutureBalancePosition struct {
	FutureBalances  Balances        `json:"balance"`
	FuturePositions FuturePositions `json:"position"`
}
