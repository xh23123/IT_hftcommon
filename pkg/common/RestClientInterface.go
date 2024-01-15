package common

import "math/big"

type RestClientInterface interface {
	BinanceRestClientInterface
	GetPremiumIndex(symbol SymbolID) []*PremiumIndexInfo
	SetMultiAssetMargin(MultiAssetMargin bool)
	GetOrder(symbol SymbolID, transactionId TransactionID, origClientOrderID string) *Order
	GetOrders(symbol SymbolID, transactionId TransactionID) []*Order
	GetSpotBalance() (SpotBalance, error)
	GetMarginBalance() (MarginBalance, error)
	GetFutureBalancePosition() (WsFutureBalance, WsFuturePosition, error)
	GetCoinFutureBalancePosition() (WsFutureBalance, WsFuturePosition, error)
	GetSpotKlines(symbol SymbolID, interval IntervalID, limit int, startTime int64, endTime int64) ([]*Kline, error)
	GetFutureKlines(symbol SymbolID, interval IntervalID, limit int, startTime int64, endTime int64) ([]*Kline, error)
	GetSuggestGasPrice() (*big.Int, error)
}
