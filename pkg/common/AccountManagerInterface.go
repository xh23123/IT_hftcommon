package common

import (
	cmap "github.com/orcaman/concurrent-map"
)

type TradingInterface interface {
	WsUpdateOrderOnOrder(OrderTradeUpdateInfo)
	WsUpdateOrderOnTrade(OrderTradeUpdateInfo)
	OnError(ErrorMsg)
}
type SpotInterface interface {
	WsUpdateSpotBalance(balance SpotBalance)
}

type FutureInterface interface {
	WsUpdateFutureBalancePosition(balancePosition WsFutureBalancePosition)
}
type CoinFutureInterface interface {
	WsUpdateCoinFutureBalancePosition(balancePosition WsFutureBalancePosition)
}

type ProcessInterface interface {
	Process(event ActionEvent)
}

type UsageInterface interface {
	GetOrders(symbol SymbolID, transactionId TransactionID) []*Order
	GetAllOrders(transactionId TransactionID) cmap.ConcurrentMap
	GetBalance(asset string, transactionId TransactionID) *Balance
	GetFuturePosition(symbol SymbolID, transactionId TransactionID) *FuturePosition
}

type AccountManagerInterface interface {
	TradingInterface
	SpotInterface
	FutureInterface
	CoinFutureInterface
	ProcessInterface
	UsageInterface
	RegisterSystemSymbols(symbols []SymbolID)
}
