package common

type StrategyInterface interface {
	InitPara()
	InitVar()
	InitMyStrategy() []ActionEvent
	OnBookTick(event BookTickWs) []ActionEvent
	OnDepth(event DepthWs) []ActionEvent
	OnTick(event TickWs) []ActionEvent
	OnTrade(event OrderTradeUpdateInfo) []ActionEvent
	OnOrder(event OrderTradeUpdateInfo) []ActionEvent
	OnKlineWs(event KlineWs) []ActionEvent
	OnFutureBookTick(event BookTickWs) []ActionEvent
	OnFutureDepth(event DepthWs) []ActionEvent
	OnFutureTick(event TickWs) []ActionEvent
	OnMarkPrice(event MarkPriceWs) []ActionEvent
	OnFutureKlineWs(event KlineWs) []ActionEvent
	OnFutureAggTrade(event AggTradeWs) []ActionEvent
	OnError(event ErrorMsg) []ActionEvent
	OnTimer() []ActionEvent
	OnExit()
}
