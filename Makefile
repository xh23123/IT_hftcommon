all:
	go build -buildmode=plugin -o pluginStrategy.so pkg/examplePluginStrategy/pluginStrategy.go
	go build -buildmode=plugin -o pluginExchange.so pkg/exampleExchangeGateway/BinaInitiator.go

unit-test:
	go test -timeout 30s -v -count=1 -run '^TestUnit.*' ./...