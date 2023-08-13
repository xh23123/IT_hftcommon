all:
	go build -buildmode=plugin -o pluginStrategy.so pkg/examplePluginStrategy/pluginStrategy.go
	go build -buildmode=plugin -o pluginExchange.so pkg/exampleExchangeGateway/BinaInitiator.go
