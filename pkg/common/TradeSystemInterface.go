package common

type OrderAgent interface {
	CreateLimitBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitSpotOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerSpotOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CancelOrderByCid(exid ExchangeID, accountIndex AccountIdx, clientOrderId string, symbol string) (ActionEvent, error)
	CancelFutureOrderByCid(exid ExchangeID, accountIndex AccountIdx, clientOrderId string, symbol string) (ActionEvent, error)
	CancelAllOrder(exid ExchangeID, accountIndex AccountIdx, symbol string) ActionEvent
	CancelAllFutureOrder(exid ExchangeID, accountIndex AccountIdx, symbol string) ActionEvent

	GetOrders(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID) []*Order
	GetOrder(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID, clientOrderId string) *Order

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
	GetSpotBalance(exid ExchangeID, accountIndex AccountIdx, asset string) *Balance
	GetFutureBalance(exid ExchangeID, accountIndex AccountIdx, asset string) *Balance
	GetFuturePosition(exid ExchangeID, accountIndex AccountIdx, symbol string) *FuturePosition
	SetMultiAssetMargin(exid ExchangeID, accountIndex AccountIdx, MultiAssetMargin bool) ActionEvent
}

type MarketAgent interface {
	InitMdConfig(*StrategyCfg)
	ResetMarketWs(exid ExchangeID, data []ResetID) ActionEvent
}

type SystemAgent interface {
	GenOrderClientId(exid ExchangeID, accountIndex AccountIdx, dataId DataID, sequence int64) string
}

type TradeSystemAgent interface {
	OrderAgent
	TimeStampAgent
	AccountAgent
	MarketAgent
	SystemAgent
	NewRestClient(exid ExchangeID, config map[string]string) RestClientInterface
	RegisterSymbols(symbols []string)
}
