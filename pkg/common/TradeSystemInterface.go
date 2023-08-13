package common

import "gopkg.in/ini.v1"

type OrderAgent interface {
	CreateLimitBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitSpotOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerSpotOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerBothFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitMakerLongFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerShortFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CreateLimitMarginOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitMakerMarginOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent

	CreateLimitBothCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64) ActionEvent
	CreateLimitLongCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent
	CreateLimitShortCoinFutureOrder(exid ExchangeID, accountIndex AccountIdx, cid string, symbol string, size float64, price float64, reduceOnly bool) ActionEvent

	CancelOrderByCid(exid ExchangeID, accountIndex AccountIdx, clientOrderId string, symbol string, transactionId TransactionID) (ActionEvent, error)
	CancelAllOrders(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID) ActionEvent

	GetOrders(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID) []*Order
	GetOrder(exid ExchangeID, accountIndex AccountIdx, symbol string, transactionId TransactionID, clientOrderId string) *Order

	ActionProcess(actions []ActionEvent)
}

type OrderFeedbackInterface interface {
	OnTrade(event OrderTradeUpdateInfo) []ActionEvent
	OnOrder(event OrderTradeUpdateInfo) []ActionEvent
	OnError(event ErrorMsg) []ActionEvent
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
	GetQueue(symbol string) chan DataEvent
	EnQueue(symbol string, event *DataEvent)
}

type MarketDataAgent interface {
	InitMdConfig(*StrategyCfg)
	ResetMarketWs(exid ExchangeID, data []ResetID) ActionEvent
	StrategyManagerCfg() StrategyManagerCfg
}

type SystemAgent interface {
	GenOrderClientId(exid ExchangeID, accountIndex AccountIdx, dataId DataID, sequence int64) string
}

type TradeSystemAgent interface {
	OrderAgent
	OrderFeedbackInterface
	TimeStampAgent
	AccountAgent
	MarketDataAgent
	SystemAgent
	GatewayInterface
	Config() *ini.File
	NewRestClient(exid ExchangeID, config map[string]string) RestClientInterface
	NewOrderManager(ExchangeID, AccountIdx, TransactionID) OrderManagerInterface
	NewBalanceManager(ExchangeID, AccountIdx) BalanceManagerInterface
	RegisterAccountManager(ExchangeID, AccountManagerInterface)
	RegisterAccountWs(ExchangeID, AccountIdx, AccountWsInterface)
	RegisterMarketWs(ExchangeID, AccountIdx, MarketWsInterface)
	RegisterSymbols(symbols []string)
}
