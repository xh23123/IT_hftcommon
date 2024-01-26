package common

type BalanceManagerInterface interface {
	InitBalances(transactionId TransactionID, balances Balances)               //should be called at the beginning of the program
	InitFuturePosition(transactionId TransactionID, positions FuturePositions) //should be called at the beginning of the program
	//spot and future
	WsUpdateBalanceInterface
	BalanceUserInterface

	//future only
	WsUpdateFuturePositionInterface
	FuturePositionUserInterface
}
