package common

import cmap "github.com/orcaman/concurrent-map"

type OrderManagerInterface interface {
	TradingInterface

	OpenOrder(orderId string) *Order
	OpenOrderByCid(orderClientId string) *Order
	OpenOrdersBySymbol(symbol SymbolID) []*Order
	AllOpenOrders() cmap.ConcurrentMap
	CreateOrderProcess(event interface{}, handler func(data *Order) (id string, err error))
	AmendOrderProcess(event interface{}, handler func(data *Order) (id string, err error))
	CancelOrderProcess(event interface{}, handler func(data CancelInfo) error)
	CancelAllOrderProcess(event interface{}, handler func(data CancelInfo))
	SetOpenOrder(orders []*Order)
	UpdateOpenOrder(openOrders []*Order) error
}
