package common

import cmap "github.com/orcaman/concurrent-map"

type BalanceManagerInterface interface {
	//spot and future
	WsUpdateBalanceInterface
	BalanceUserInterface
	SetBalances(transactionId TransactionID, balances cmap.ConcurrentMap)

	//future only
	WsUpdateFuturePositionInterface
	FuturePositionUserInterface
	SetFuturePosition(transactionId TransactionID, positions cmap.ConcurrentMap)
}
