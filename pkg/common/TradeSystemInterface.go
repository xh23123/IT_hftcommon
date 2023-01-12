package common

type OrderAgent interface {
	CreateLimitBothFutureOrder(exid ExchangeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitLongFutureOrder(exid ExchangeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitShortFutureOrder(exid ExchangeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitSpotOrder(exid ExchangeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerSpotOrder(exid ExchangeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerBothFutureOrder(exid ExchangeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerLongFutureOrder(exid ExchangeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerShortFutureOrder(exid ExchangeID, cid string, symbol string, size float64, price float64) ActionEvent

	CancelOrderByCid(exid ExchangeID, clientOrderId string, symbol string) (ActionEvent, error)
	CancelFutureOrderByCid(exid ExchangeID, clientOrderId string, symbol string) (ActionEvent, error)

	GetSpotOrders(exid ExchangeID, symbol string) []*Order
	GetSpotOrder(exid ExchangeID, symbol string, clientOrderId string) *Order
	GetFutureOrders(exid ExchangeID, symbol string) []*Order
	GetFutureOrder(exid ExchangeID, symbol string, clientOrderId string) *Order

	ActionProcess(actions []ActionEvent)
}

type TimeStampAgent interface {
	OrderTimestamp(dataExid ExchangeID,
		dataId DataID,
		dataTimestamp int64,
		orderAction *ActionEvent,
	)
}

type AccountAgent interface {
	GetSpotBalance(exid ExchangeID, asset string) *Balance
	GetFutureBalance(exid ExchangeID, asset string) *Balance
	GetFuturePosition(exid ExchangeID, symbol string) *FuturePosition
	SetMultiAssetMargin(exid ExchangeID, MultiAssetMargin bool) ActionEvent
}

type MarketAgent interface {
	InitMdConfig(*StrategyCfg)
}

type SystemAgent interface {
	GenOrderClientId(dataId DataID, exchangeID ExchangeID, sequence int64)
}

type TradeSystemAgent interface {
	OrderAgent
	TimeStampAgent
	AccountAgent
	MarketAgent
	SystemAgent
	NewRestClient(exid ExchangeID, config map[string]string) RestClientInterface
}
