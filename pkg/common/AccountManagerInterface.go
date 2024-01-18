package common

import (
	cmap "github.com/orcaman/concurrent-map"
)

type TradingInterface interface {
	WsUpdateOrderOnOrder(*OrderTradeUpdateInfo)
	WsUpdateOrderOnTrade(*OrderTradeUpdateInfo)
	OnError(*ErrorMsg)
}

type WsUpdateBalanceInterface interface {
	//spot and future
	WsUpdateBalance(transactionId TransactionID, balance cmap.ConcurrentMap)
}

type BalanceUserInterface interface {
	GetBalance(asset SymbolID, transactionId TransactionID) *Balance
}

type WsUpdateFuturePositionInterface interface {
	WsUpdateFuturePosition(transactionId TransactionID, positions FuturePositions)
}

type FuturePositionUserInterface interface {
	GetFuturePosition(symbol SymbolID, transactionId TransactionID) *FuturePosition
}

type ProcessInterface interface {
	Process(event *ActionEvent)
}

type UsageInterface interface {
	BalanceUserInterface
	FuturePositionUserInterface
	GetOrders(symbol SymbolID, transactionId TransactionID) []*Order
	GetAllOrders(transactionId TransactionID) cmap.ConcurrentMap
}

type AccountManagerInterface interface {
	TradingInterface
	WsUpdateBalanceInterface
	WsUpdateFuturePositionInterface
	ProcessInterface
	UsageInterface
	RegisterSystemSymbols(symbols []SymbolID)
}
