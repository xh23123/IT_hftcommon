all:
	go build -buildmode=plugin -o pluginStrategy.so pkg/examplePluginStrategy/pluginStrategy.go

unit-test:
	go test -timeout 30s -v -count=1 -run '^TestUnit.*' ./...
