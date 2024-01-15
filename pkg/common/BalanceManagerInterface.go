package common

import cmap "github.com/orcaman/concurrent-map"

type BalanceManagerInterface interface {
	//spot and future
	SetBalances(transactionId TransactionID, balances cmap.ConcurrentMap)
	WsUpdateBalance(transactionId TransactionID, balance *Balance)
	GetBalance(transactionId TransactionID, asset string) *Balance

	//future only
	SetFuturePosition(transactionId TransactionID, position cmap.ConcurrentMap)
	GetFuturePosition(symbol SymbolID, transactionId TransactionID) *FuturePosition
	WsUpdateFuturePosition(transactionId TransactionID, position WsFuturePosition)
}
