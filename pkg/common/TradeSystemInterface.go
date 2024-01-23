package common

import "gopkg.in/ini.v1"

type OrderOptions struct {
	ReduceOnly bool       `json:"reduce_only"` //For future BOTH orders only
	PositionId PositionID `json:"position_id"` //For future orders
}

type OrderAgent interface {
	CreateNewOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, clientOrderId ClientOrderidID, symbol SymbolID, transactionId TransactionID, size float64, price float64, orderOptions *OrderOptions) *ActionEvent
	CreateAmendOrder(exid ExchangeID, accountIndex AccountIdx, orderType OrderTypeID, clientOrderId ClientOrderidID, symbol SymbolID, transactionId TransactionID, size float64, price float64, orderOptions *OrderOptions) *ActionEvent
	CancelOrderByCid(exid ExchangeID, accountIndex AccountIdx, clientOrderId ClientOrderidID, symbol SymbolID, transactionId TransactionID) (*ActionEvent, error)
	CancelAllOrders(exid ExchangeID, accountIndex AccountIdx, symbol SymbolID, transactionId TransactionID) *ActionEvent

	GetOrders(exid ExchangeID, accountIndex AccountIdx, symbol SymbolID, transactionId TransactionID) []*Order
	GetOrder(exid ExchangeID, accountIndex AccountIdx, symbol SymbolID, transactionId TransactionID, clientOrderId ClientOrderidID) *Order

	ActionProcess(actions []*ActionEvent)
}

type OrderFeedbackInterface interface {
	OnError(event *ErrorMsg)
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
		symbol SymbolID,
		actionType ActionID,
		timeEventID TimeEventID,
		clientOrderId ClientOrderidID,
		orderId OrderidID,
		feedBackTimestamp int64,
	)
}

type AccountAgent interface {
	GetBalance(exid ExchangeID, accountIndex AccountIdx, asset SymbolID, transactionId TransactionID) *Balance
	GetFuturePosition(exid ExchangeID, accountIndex AccountIdx, symbol SymbolID, transactionId TransactionID) *FuturePosition
	SetMultiAssetMargin(exid ExchangeID, accountIndex AccountIdx, MultiAssetMargin bool) *ActionEvent
	SetDualSidePosition(exid ExchangeID, accountIndex AccountIdx, transactionId TransactionID, dualSidePosition bool) *ActionEvent

	WsUpdateBalance(exid ExchangeID, accountIndex AccountIdx, transactionId TransactionID, balances Balances)
	WsUpdateFuturePosition(exid ExchangeID, accountIndex AccountIdx, transactionId TransactionID, positions FuturePositions)
}

type GatewayInterface interface {
	StartGateWay()
	EnQueue(symbol SymbolID, event *DataEvent)
}

type MarketDataAgent interface {
	InitMdConfig(*StrategyCfg)
	ResetMarketWs(exid ExchangeID, data []ResetID) *ActionEvent
	MarketDataConfigs() MarketDataConfigs
}

type SystemAgent interface {
	GenOrderClientId(exid ExchangeID, accountIndex AccountIdx, dataId DataID, sequence int64) ClientOrderidID
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
	RegisterSymbols(symbols []SymbolID)

	//options: {"marginMode": "unimargin"}  -> unimargin mode
	//options: {"marginMode": "normal"}  -> normal mode
	AmendType(exid ExchangeID, transactionId TransactionID, options map[string]string) AmendTypeID
}
