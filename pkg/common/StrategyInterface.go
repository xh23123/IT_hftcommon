package common

type StrategyInterface interface {
	InitPara()
	InitVar()
	InitMyStrategy() []ActionEvent
	OnBookTick(event BookTickWs) []ActionEvent
	OnDepth(event DepthWs) []ActionEvent
	OnTick(event TickWs) []ActionEvent
	OnTrade(event TradeUpdateInfo) []ActionEvent
	OnOrder(event OrderUpdateInfo) []ActionEvent
	OnKlineWs(event KlineWs) []ActionEvent
	OnFutureBookTick(event BookTickWs) []ActionEvent
	OnFutureDepth(event DepthWs) []ActionEvent
	OnFutureTick(event TickWs) []ActionEvent
	OnMarkPrice(event MarkPriceWs) []ActionEvent
	OnFutureKlineWs(event KlineWs) []ActionEvent
	OnFutureAggTrade(event AggTradeWs) []ActionEvent
	OnFutureTrade(event TradeUpdateInfo) []ActionEvent
	OnFutureOrder(event OrderUpdateInfo) []ActionEvent
	OnError(event ErrorMsg) []ActionEvent
	OnTimer() []ActionEvent
	OnExit()
}
