package common

type BalanceManagerInterface interface {
	//spot and future
	WsUpdateBalanceInterface
	BalanceUserInterface
	SetBalances(transactionId TransactionID, balances Balances)

	//future only
	WsUpdateFuturePositionInterface
	FuturePositionUserInterface
	SetFuturePosition(transactionId TransactionID, positions FuturePositions)
}
