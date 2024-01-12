package common

type OrderEntryCallback interface {
	OnOrder(event OrderTradeUpdateInfo) []ActionEvent
	OnTrade(event OrderTradeUpdateInfo) []ActionEvent
	OnError(event ErrorMsg) []ActionEvent
}

type MarketDataCallback interface {
	OnBookTick(event BookTickWs) []ActionEvent
	OnDepth(event DepthWs) []ActionEvent
	OnTick(event TickWs) []ActionEvent
	OnKlineWs(event KlineWs) []ActionEvent
	OnTradeWs(event TradeWs) []ActionEvent
	OnMarkPrice(event MarkPriceWs) []ActionEvent
	OnOrderbook(event Orderbook) []ActionEvent
	OnDexBookTicks(event DexBookTicks) []ActionEvent
	OnDexTrades(event DexTrades) []ActionEvent
}

type StrategyCommonFunctions interface {
	InitPara()
	InitVar()
	InitMyStrategy() []ActionEvent
	OnTimer() []ActionEvent
	OnExit()
}

type StrategyInterface interface {
	OrderEntryCallback
	MarketDataCallback
	StrategyCommonFunctions
}
