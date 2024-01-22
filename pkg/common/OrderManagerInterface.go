package common

import cmap "github.com/orcaman/concurrent-map"

type OrderManagerInterface interface {
	TradingInterface

	OpenOrder(orderId OrderidID) *Order
	OpenOrderByCid(orderClientId OrderClientidID) *Order
	OpenOrdersBySymbol(symbol SymbolID) []*Order
	AllOpenOrders() cmap.ConcurrentMap
	CreateOrderProcess(event interface{}, handler func(data *Order) (id OrderidID, err error))
	AmendOrderProcess(event interface{}, handler func(data *Order) (id OrderidID, err error))
	CancelOrderProcess(event interface{}, handler func(data CancelInfo) error)
	CancelAllOrderProcess(event interface{}, handler func(data CancelInfo))
	SetOpenOrder(orders []*Order)
	UpdateOpenOrder(openOrders []*Order) error
}
