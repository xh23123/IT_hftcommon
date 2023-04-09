package common

type OrderAgent interface {
	CreateLimitBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitSpotOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerSpotOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitMakerShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent

	CreateLimitMarginOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CreateLimitBothCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitLongCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitShortCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent

	CancelOrderByCid(exid ExchangeID, accountIndex AccountIdx, clientOrderId string, symbol string, transactionId TransactionID) (ActionEvent, error)
	CancelAllOrders(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID) ActionEvent

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
	GetBalance(exid ExchangeID, accountIndex AccountIdx, asset string, transactionId TransactionID) *Balance
	GetFuturePosition(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID) *FuturePosition
	SetMultiAssetMargin(exid ExchangeID, accountIndex AccountIdx, transactionId TransactionID, MultiAssetMargin bool) ActionEvent
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
