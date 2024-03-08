package common

type AccountManagerInterface interface {
	Process(event *ActionEvent)
	RegisterSystemSymbols(symbols []SymbolID)
}
