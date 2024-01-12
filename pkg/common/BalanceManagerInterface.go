package common

type BalanceManagerInterface interface {
	InitBalances(transactionId TransactionID, balances Balances)                //should be called at the beginning of the program
	InitFuturePositions(transactionId TransactionID, positions FuturePositions) //should be called at the beginning of the program
}
