all:
	go build -buildmode=plugin -o pluginStrategy.so pkg/examplePluginStrategy/pluginStrategy.go
