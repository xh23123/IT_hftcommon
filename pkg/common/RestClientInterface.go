package common

import "math/big"

type RestClientInterface interface {
	BinanceRestClientInterface
	GetPremiumIndex(symbol SymbolID) []*PremiumIndexInfo
	SetMultiAssetMargin(MultiAssetMargin bool)
	GetOrder(symbol SymbolID, transactionId TransactionID, orderClientId string) *Order
	GetOrders(symbol SymbolID, transactionId TransactionID) []*Order
	GetBalances(transactionId TransactionID) (Balances, error)
	GetMarginBalances() (MarginBalances, error)
	GetFuturePositions(transactionId TransactionID) (FuturePositions, error)
	GetKlines(symbol SymbolID, transactionId TransactionID, interval IntervalID, limit int, startTime int64, endTime int64) ([]*Kline, error)
	GetSuggestGasPrice() (*big.Int, error)
}
