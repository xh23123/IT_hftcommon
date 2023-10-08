package common

import cmap "github.com/orcaman/concurrent-map"

type OrderManagerInterface interface {
	OpenOrder(orderId string) (*Order, bool)           //TODO
	OpenOrdersBySymbol(symbol string) ([]*Order, bool) //TODO
	AllOpenOrders() cmap.ConcurrentMap                 //TODO
	CreateOrderProcess(event interface{}, handler func(data *Order) (id string, err error))
	CancelOrderProcess(event interface{}, handler func(data CancelInfo) error)
	CancelAllOrderProcess(event interface{}, handler func(data CancelInfo))
	WsUpdateOrderOnOrder(OrderTradeUpdateInfo)
	SetOpenOrder(orders []*Order)
	WsUpdateOrderOnTrade(event OrderTradeUpdateInfo)
	UpdateOpenOrder(openOrders []*Order) error
}
