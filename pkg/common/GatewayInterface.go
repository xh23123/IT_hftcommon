package common

type GatewayInterface interface {
	StartGateWay()
	EnQueue(event *DataEvent)
}
