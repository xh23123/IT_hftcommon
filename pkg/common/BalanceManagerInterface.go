package common

import cmap "github.com/orcaman/concurrent-map"

type BalanceManagerInterface interface {
	//spot and future
	WsBalanceInterface
	BalanceUserInterface
	SetBalances(transactionId TransactionID, balances cmap.ConcurrentMap)

	//future only
	WsFuturePositionInterface
	FuturePositionUserInterface
	SetFuturePosition(transactionId TransactionID, position cmap.ConcurrentMap)
}
