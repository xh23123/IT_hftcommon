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
	OnFutureBookTick(event BookTickWs) []ActionEvent
	OnFutureDepth(event DepthWs) []ActionEvent
	OnFutureTick(event TickWs) []ActionEvent
	OnMarkPrice(event MarkPriceWs) []ActionEvent
	OnFutureKlineWs(event KlineWs) []ActionEvent
	OnFutureAggTrade(event TradeWs) []ActionEvent
	OnCoinFutureBookTick(event BookTickWs) []ActionEvent
	OnDexBookTicks(event DexBookTicks) []ActionEvent
	OnDexTrades(event DexTrades) []ActionEvent
	OnFutureOrderbook(event Orderbook) []ActionEvent
}

type StrategyInterface interface {
	OrderEntryCallback
	MarketDataCallback

	InitPara()
	InitVar()
	InitMyStrategy() []ActionEvent
	OnTimer() []ActionEvent
	OnExit()
}
