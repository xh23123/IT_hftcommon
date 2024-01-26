package common

type BalanceManagerInterface interface {
	InitBalances(transactionId TransactionID, balances Balances)
	InitFuturePosition(transactionId TransactionID, positions FuturePositions)
	//spot and future
	WsUpdateBalanceInterface
	BalanceUserInterface

	//future only
	WsUpdateFuturePositionInterface
	FuturePositionUserInterface
}
