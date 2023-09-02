package common

type TimeEventID string

var (
	DataArriveToGateway   TimeEventID = "1"
	DataArriveToStrategy  TimeEventID = "2"
	OrderSendFromStrategy TimeEventID = "3"
	OrderArriveToGateway  TimeEventID = "4"
	OrderSendFromGateway  TimeEventID = "5"
	OnOrderToGateway      TimeEventID = "6"
	OnOrderToStrategy     TimeEventID = "7"
	OnTradeToGateway      TimeEventID = "8"
	OnTradeToStrategy     TimeEventID = "9"
	OnCancelToGateway     TimeEventID = "10"
	OnCancelToStrategy    TimeEventID = "11"
)
