package common

type RestClientInterface interface {
	BinanceRestClientInterface
	GetPremiumIndex(symbol string) []*PremiumIndexInfo
	SetMultiAssetMargin(MultiAssetMargin bool)
	GetOrder(symbol string, transactionId TransactionID, origClientOrderID string) *Order
	GetOrders(symbol string, transactionId TransactionID) []*Order
	GetSpotBalance() (SpotBalance, error)
	GetMarginBalance() (MarginBalance, error)
	GetFutureBalancePosition() (WsFutureBalance, WsFuturePosition, error)
	GetCoinFutureBalancePosition() (WsFutureBalance, WsFuturePosition, error)
	GetSpotKlines(symbol string, interval IntervalID, limit int, startTime int64, endTime int64) ([]*Kline, error)
	GetFutureKlines(symbol string, interval IntervalID, limit int, startTime int64, endTime int64) ([]*Kline, error)
}
