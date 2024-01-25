package common

type OrderManagerInterface interface {
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
