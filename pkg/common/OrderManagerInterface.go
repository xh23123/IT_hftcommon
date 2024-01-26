package common

type OrderManagerInterface interface {
	IniOpenOrder(orders []*Order) //should be called at the beginning of the program
	TradingInterface

	OpenOrder(orderId OrderidID) *Order
	OpenOrderByCid(orderClientId ClientOrderidID) *Order
	OpenOrdersBySymbol(symbol SymbolID) []*Order
	AllOpenOrders() []*Order
	CreateOrderProcess(event interface{}, handler func(data *Order))
	AmendOrderProcess(event interface{}, handler func(data *Order))
	CancelOrderProcess(event interface{}, handler func(data *CancelInfo))
	CancelAllOrderProcess(event interface{}, handler func(data *CancelInfo))
}
