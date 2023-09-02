package common

import cmap "github.com/orcaman/concurrent-map"

type BalanceManagerInterface interface {
	SetSpotBalance(balances cmap.ConcurrentMap)
	SetFutureBalancePosition(balance cmap.ConcurrentMap, position cmap.ConcurrentMap)
	WsUpdateSpotBalance(balance SpotBalance)
	GetSpotBalance(asset string) *Balance
	GetFutureBalance(asset string) *Balance
	GetFuturePosition(symbol string, transactionId TransactionID) *FuturePosition
	WsUpdateCoinFutureBalancePosition(balance WsFutureBalance, position WsFuturePosition)
	WsUpdateFutureBalancePosition(balance WsFutureBalance, position WsFuturePosition)
}
