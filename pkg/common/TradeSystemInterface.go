package common

import "gopkg.in/ini.v1"

type OrderAgent interface {

	//GFD /LIMIT orders
	CreateLimitSpotOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerSpotOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CreateLimitBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CreateLimitMakerBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitMakerLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CreateLimitMarginOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerMarginOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CreateLimitBothCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitLongCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitShortCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	//IOC orders
	CreateIocSpotOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateIocMarginOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CreateIocBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateIocLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateIocShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CreateIocBothCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateIocLongCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateIocShortCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	//Amend orders
	CreateAmendSpotOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateAmendMarginOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, cid string, symbol string, size float64, price float64) ActionEvent

	CreateAmendBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateAmendLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateAmendShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, cid string, symbol string, size float64, price float64) ActionEvent

	CreateAmendBothCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateAmendLongCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, cid string, symbol string, size float64, price float64) ActionEvent
	CreateAmendShortCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, cid string, symbol string, size float64, price float64) ActionEvent

	CancelOrderByCid(exid ExchangeID, accountIndex AccountIdx, clientOrderId string, symbol string, transactionId TransactionID) (ActionEvent, error)
	CancelAllOrders(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID) ActionEvent

	GetOrders(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID) []*Order
	GetOrder(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID, clientOrderId string) *Order

	ActionProcess(actions []ActionEvent)
}

type OrderFeedbackInterface interface {
	OnError(event ErrorMsg)
}

type TimeStampAgent interface {
	OrderTimestamp(dataExid ExchangeID,
		dataId DataID,
		dataTimestamp int64,
		orderAction *ActionEvent,
	)

	FeedbackTimestamp(dataExid ExchangeID,
		accountIndex AccountIdx,
		dataId DataID,
		symbol string,
		actionType ActionID,
		timeEventID TimeEventID,
		orderClientId string,
		orderId string,
		feedBackTimestamp int64,
	)
}

type AccountAgent interface {
	GetBalance(exid ExchangeID, accountIndex AccountIdx, asset string, transactionId TransactionID) *Balance
	GetFuturePosition(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID) *FuturePosition
	SetMultiAssetMargin(exid ExchangeID, accountIndex AccountIdx, MultiAssetMargin bool) ActionEvent
	SetDualSidePosition(exid ExchangeID, accountIndex AccountIdx, transactionId TransactionID, dualSidePosition bool) ActionEvent

	WsUpdateFutureBalancePosition(exid ExchangeID, accountIndex AccountIdx, balancePosition WsFutureBalancePosition)
}

type GatewayInterface interface {
	StartGateWay()
	EnQueue(symbol string, event *DataEvent)
}

type MarketDataAgent interface {
	InitMdConfig(*StrategyCfg)
	ResetMarketWs(exid ExchangeID, data []ResetID) ActionEvent
	MarketDataConfigs() MarketDataConfigs
}

type SystemAgent interface {
	GenOrderClientId(exid ExchangeID, accountIndex AccountIdx, dataId DataID, sequence int64) string
}

type DebugInterface interface {
	SetOpenOrder(exid ExchangeID, accountIndex AccountIdx, transactionID TransactionID, orders []*Order)
}

type TradeSystemAgent interface {
	OrderAgent
	OrderFeedbackInterface
	TimeStampAgent
	AccountAgent
	MarketDataAgent
	SystemAgent
	GatewayInterface
	DebugInterface
	Config() *ini.File
	NewRestClient(exid ExchangeID, config map[string]string) RestClientInterface
	NewOrderManager(ExchangeID, AccountIdx, TransactionID) OrderManagerInterface
	NewBalanceManager(ExchangeID, AccountIdx) BalanceManagerInterface
	RegisterAccountManager(ExchangeID, AccountManagerInterface)
	RegisterAccountWs(ExchangeID, AccountIdx, AccountWsInterface)
	RegisterMarketWs(ExchangeID, AccountIdx, MarketWsInterface)
	RegisterSymbols(symbols []string)

	//options: {"marginMode": "unimargin"}  -> unimargin mode
	//options: {"marginMode": "normal"}  -> normal mode
	AmendType(exid ExchangeID, transactionId TransactionID, options map[string]string) AmendTypeID
}
