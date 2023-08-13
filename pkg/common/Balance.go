package common

type WsFutureBalancePosition struct {
	FutureBalances  WsFutureBalance  `json:"balance"`
	FuturePositions WsFuturePosition `json:"position"`
}
