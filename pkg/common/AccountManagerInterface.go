package common

import (
	cmap "github.com/orcaman/concurrent-map"
)

type TradingInterface interface {
	WsUpdateOrderOnOrder(OrderTradeUpdateInfo)
	WsUpdateOrderOnTrade(OrderTradeUpdateInfo)
	OnError(ErrorMsg)
}

type WsBalanceInterface interface {
	//spot and future
	WsUpdateBalance(transactionId TransactionID, balance *Balance)
}

type BalanceUserInterface interface {
	GetBalance(asset string, transactionId TransactionID) *Balance
}

type WsFuturePositionInterface interface {
	WsUpdateFuturePosition(transactionId TransactionID, position WsFuturePosition)
}

type FuturePositionUserInterface interface {
	GetFuturePosition(symbol SymbolID, transactionId TransactionID) *FuturePosition
}

type ProcessInterface interface {
	Process(event ActionEvent)
}

type UsageInterface interface {
	BalanceUserInterface
	FuturePositionUserInterface
	GetOrders(symbol SymbolID, transactionId TransactionID) []*Order
	GetAllOrders(transactionId TransactionID) cmap.ConcurrentMap
}

type AccountManagerInterface interface {
	TradingInterface
	WsBalanceInterface
	WsFuturePositionInterface
	ProcessInterface
	UsageInterface
	RegisterSystemSymbols(symbols []SymbolID)
}
