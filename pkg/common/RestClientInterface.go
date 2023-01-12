package common

type RestClientInterface interface {
	GetPremiumIndex(symbol string) []*PremiumIndexInfo
	SetMultiAssetMargin(MultiAssetMargin bool)
	GetSpotOrder(symbol string, origClientOrderID string) *Order
	GetSpotOrders(symbol string) []*Order
	GetFutureOrder(symbol string, origClientOrderID string) *Order
	GetFutureOrders(symbol string) []*Order
	GetSpotBalance() WsSpotBalance
	GetFutureBalancePosition() (WsFutureBalance, WsFuturePosition, error)
	GetSpotKlines(symbol string, interval IntervalID, limit int, startTime int64, endTime int64) ([]*Kline, error)
	GetFutureKlines(symbol string, interval IntervalID, limit int, startTime int64, endTime int64) ([]*Kline, error)
}
