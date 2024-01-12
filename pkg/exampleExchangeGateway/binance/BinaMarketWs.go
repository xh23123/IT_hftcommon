package bina

import (
	"fmt"
	"time"

	"github.com/xh23123/IT_hftcommon/pkg/common"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2/futures"
	"go.uber.org/zap"
)

var _ common.MarketWsInterface = (*BinaMarketWs)(nil)
var notify_interval int64 = 500 * int64(time.Millisecond)

const MarketWsAccountIndex = 0
const BinanceMdWsAccountIdx = 0

type BinaMarketWs struct {
	systemAgent    common.TradeSystemAgent
	name           string
	tickTrigRecord *ConcurrentTickTriggerInfoMap
	cancelCtx      *ConcurrentCancelMap
}

func NewMarketWs(systemAgent common.TradeSystemAgent, info AccountInfo) {

	if info.Index != MarketWsAccountIndex {
		return
	}

	bmw := &BinaMarketWs{
		systemAgent: systemAgent,
		name:        "BinaMarketWs",
		tickTrigRecord: &ConcurrentTickTriggerInfoMap{
			TickTrigRecord: make(map[common.TransactionID]map[string]*common.TickTrigInfo),
		},
		cancelCtx: &ConcurrentCancelMap{
			CancelCtx: make(map[common.ResetID]chan struct{}),
		},
	}

	systemAgent.RegisterMarketWs(common.BINANCEID, MarketWsAccountIndex, bmw)
	bmw.registerMarketWs()
}

func (b *BinaMarketWs) registerMarketWs() {
	b.systemAgent.StartGateWay()
	go b.registerSpotBookTick()
	go b.registerSpotTrade()
	go b.registerFutureBookTick()
	go b.registerFutureDepth()
	go b.registerFutureAggTrade()
}
func (b *BinaMarketWs) errHandler(err error) {
	common.Logger.Error(b.name+" BinaMarketWs : ",
		zap.String("error", err.Error()))
}

func (b *BinaMarketWs) shouldNotifyBookTicker(transactionId common.TransactionID, symbol string, data *common.BookTickWs) bool {
	if prev, ok := b.tickTrigRecord.GetTriggerInfo(transactionId, symbol); ok {
		intervalinfo := b.systemAgent.StrategyManagerCfg().EXWsCfg[common.BINANCEID].IntervalMap[symbol]
		if (prev.PrevBap != data.BestAskPrice || prev.PrevBbp != data.BestBidPrice) &&
			data.ReceiveTimestamp-prev.PrevTime >= intervalinfo.TickTrigInterval ||
			(data.ReceiveTimestamp-prev.PrevTime >= notify_interval) {

			prev.PrevTime = data.ReceiveTimestamp
			prev.PrevBap = data.BestAskPrice
			prev.PrevBbp = data.BestBidPrice
			return true
		} else {
			return false
		}
	} else {
		prev.PrevTime = data.ReceiveTimestamp
		prev.PrevBap = data.BestAskPrice
		prev.PrevBbp = data.BestBidPrice
	}
	return true
}

func (b *BinaMarketWs) registerSpotBookTick() {

	wsHandler := func(event *binance.WsBookTickerEvent) {
		symbol := GetSysSymbol(event.Symbol)
		data := b.bookTickHandler(event)

		dataEvent := common.NewDataEvent(common.BINANCEID, BinanceMdWsAccountIdx, common.SpotID, common.BookTickID, symbol, data)

		if b.shouldNotifyBookTicker(common.SpotID, symbol, &data) {
			b.systemAgent.EnQueue(symbol, dataEvent)
		}

	}
	for {
		ri := b.systemAgent.StrategyManagerCfg().EXWsCfg[common.BINANCEID].RegisterInfo.RegisterSpotBookTick

		if len(ri) == 0 {
			stopC := make(chan struct{})
			b.cancelCtx.SetCancelChannel(common.ResetSpotBookTick, stopC)
			<-stopC
			time.Sleep(100 * time.Millisecond)
			continue
		}

		common.Logger.Info(b.name + " Establishing Spot BookTick connection. Please wait")
		doneC, stopC, err := binance.WsCombinedBookTickerServe(ConvertToExchangeSymbolsSlice(ri), wsHandler, b.errHandler)
		if err != nil {
			common.Logger.Error(b.name+" Spot BookTick Websocket Error:", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		b.cancelCtx.SetCancelChannel(common.ResetSpotBookTick, stopC)
		<-doneC
		time.Sleep(100 * time.Millisecond)
	}
}

func (b *BinaMarketWs) registerFutureBookTick() {

	wsHandler := func(event *futures.WsBookTickerEvent) {
		symbol := GetSysSymbol(event.Symbol)

		data := b.futureBookTickHandler(event)

		dataEvent := common.NewDataEvent(common.BINANCEID, BinanceMdWsAccountIdx, common.FutureID, common.BookTickID, symbol, data)
		if b.shouldNotifyBookTicker(common.SpotID, symbol, &data) {
			b.systemAgent.EnQueue(symbol, dataEvent)
		}
	}
	for {
		ri := b.systemAgent.StrategyManagerCfg().EXWsCfg[common.BINANCEID].RegisterInfo.RegisterFutureBookTick
		if len(ri) == 0 {
			stopC := make(chan struct{})
			b.cancelCtx.SetCancelChannel(common.ResetFutureBookTick, stopC)
			<-stopC
			time.Sleep(100 * time.Millisecond)
			continue
		}

		common.Logger.Info(b.name + " Establishing Future BookTick connection. Please wait")
		doneC, stopC, err := futures.WsCombinedBookTickerServe(ConvertToExchangeSymbolsSlice(ri), wsHandler, b.errHandler)
		if err != nil {
			common.Logger.Error(b.name+" Future BookTick Websocket Error:", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		b.cancelCtx.SetCancelChannel(common.ResetFutureBookTick, stopC)
		<-doneC
		time.Sleep(100 * time.Millisecond)
	}
}

func (b *BinaMarketWs) bookTickHandler(data *binance.WsBookTickerEvent) common.BookTickWs {
	askPrice := common.Str2Float(data.BestAskPrice)
	askSize := common.Str2Float(data.BestAskQty)
	bidPrice := common.Str2Float(data.BestBidPrice)
	bidSize := common.Str2Float(data.BestBidQty)
	return common.BookTickWs{
		MarketDataHeader: common.MarketDataHeader{
			Exchange:         common.BINANCEID,
			DataID:           common.BookTickID,
			Symbol:           GetSysSymbol(data.Symbol),
			ReceiveTimestamp: time.Now().UnixNano(),
		},

		BestBidPrice: bidPrice,
		BestBidSize:  bidSize,
		BestAskPrice: askPrice,
		BestAskSize:  askSize,
	}
}

func (b *BinaMarketWs) futureBookTickHandler(data *futures.WsBookTickerEvent) common.BookTickWs {
	askPrice := common.Str2Float(data.BestAskPrice)
	askSize := common.Str2Float(data.BestAskQty)
	bidPrice := common.Str2Float(data.BestBidPrice)
	bidSize := common.Str2Float(data.BestBidQty)
	return common.BookTickWs{
		MarketDataHeader: common.MarketDataHeader{
			Exchange:         common.BINANCEID,
			DataID:           common.BookTickID,
			Symbol:           GetSysSymbol(data.Symbol),
			ReceiveTimestamp: time.Now().UnixNano(),
		},
		BestBidPrice: bidPrice,
		BestBidSize:  bidSize,
		BestAskPrice: askPrice,
		BestAskSize:  askSize,
	}
}

func (b *BinaMarketWs) futureDepthHandler(event interface{}) common.DepthWs {
	data := event.(*futures.WsDepthEvent)

	symbol := GetSysSymbol(data.Symbol)
	bids := [20][2]float64{}
	for i, bid := range data.Bids {
		price := common.Str2Float(bid.Price)
		size := common.Str2Float(bid.Quantity)
		bids[i] = [2]float64{price, size}
	}
	asks := [20][2]float64{}
	for i, ask := range data.Asks {
		price := common.Str2Float(ask.Price)
		size := common.Str2Float(ask.Quantity)
		asks[i] = [2]float64{price, size}
	}

	return common.DepthWs{
		MarketDataHeader: common.MarketDataHeader{
			Exchange:         common.BINANCEID,
			DataID:           common.DepthID,
			Symbol:           symbol,
			ReceiveTimestamp: time.Now().UnixNano(),
		},
		Bids: bids,
		Asks: asks,
	}
}

func (b *BinaMarketWs) FutureAggTradeHandler(event interface{}) common.TradeWs {
	data := event.(*futures.WsAggTradeEvent)
	price := common.Str2Float(data.Price)
	size := common.Str2Float(data.Quantity)

	return common.TradeWs{
		MarketDataHeader: common.MarketDataHeader{
			Exchange:         common.BINANCEID,
			DataID:           common.AggTradeID,
			Symbol:           GetSysSymbol(data.Symbol),
			ReceiveTimestamp: time.Now().UnixNano(),
		},
		Price:      price,
		Size:       size,
		TradeCount: data.LastTradeID - data.FirstTradeID + 1,
		IsMaker:    data.Maker}
}

func (b *BinaMarketWs) SpotTradeHandler(event interface{}) common.TradeWs {
	data := event.(*binance.WsTradeEvent)
	price := common.Str2Float(data.Price)
	size := common.Str2Float(data.Quantity)

	return common.TradeWs{
		MarketDataHeader: common.MarketDataHeader{
			Exchange:          common.BINANCEID,
			DataID:            common.TradeID,
			Symbol:            GetSysSymbol(data.Symbol),
			ReceiveTimestamp:  time.Now().UnixNano(),
			ExchangeTimestamp: data.Time,
		},
		Price:      price,
		Size:       size,
		TradeCount: 1,
		IsMaker:    data.IsBuyerMaker,
	}
}

func (b *BinaMarketWs) registerFutureDepth() {

	wsHandler := func(event *futures.WsDepthEvent) {
		symbol := GetSysSymbol(event.Symbol)
		data := b.futureDepthHandler(event)
		dataEvent := common.NewDataEvent(common.BINANCEID, BinanceMdWsAccountIdx, common.FutureID, common.DepthID, symbol, data)
		b.systemAgent.EnQueue(symbol, dataEvent)
	}

	for {
		ri := b.systemAgent.StrategyManagerCfg().EXWsCfg[common.BINANCEID].RegisterInfo.RegisterFutureDepth
		if len(ri) == 0 {
			stopC := make(chan struct{})
			b.cancelCtx.SetCancelChannel(common.ResetFutureDepth, stopC)
			<-stopC
			time.Sleep(100 * time.Millisecond)
			continue
		}

		common.Logger.Info(b.name + " Establishing Future Depth connection. Please wait")
		doneC, stopC, err := futures.WsCombinedDepthServe(ConvertToExchangeSymbolsMap(ri), wsHandler, b.errHandler)
		if err != nil {
			common.Logger.Error(b.name+" Future Depth Websocket Error:", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		b.cancelCtx.SetCancelChannel(common.ResetFutureDepth, stopC)
		<-doneC
		time.Sleep(100 * time.Millisecond)
	}
}

func (b *BinaMarketWs) registerSpotTrade() {

	wsHandler := func(event *binance.WsCombinedTradeEvent) {
		symbol := GetSysSymbol(event.Data.Symbol)
		data := b.SpotTradeHandler(&event.Data)

		dataEvent := common.NewDataEvent(common.BINANCEID, BinanceMdWsAccountIdx, common.SpotID, common.TradeID, symbol, data)
		b.systemAgent.EnQueue(symbol, dataEvent)
	}

	for {

		ri := b.systemAgent.StrategyManagerCfg().EXWsCfg[common.BINANCEID].RegisterInfo.RegisterSpotTrade
		if len(ri) == 0 {
			stopC := make(chan struct{})
			b.cancelCtx.SetCancelChannel(common.ResetSpotTrade, stopC)
			<-stopC
			time.Sleep(100 * time.Millisecond)
			continue
		}

		common.Logger.Info(b.name + " Establishing Spot trade connection. Please wait")
		doneC, stopC, err := binance.WsCombinedTradeServe(ConvertToExchangeSymbolsSlice(ri), wsHandler, b.errHandler)
		if err != nil {
			common.Logger.Error(b.name+" Spot trade Websocket Error:", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		b.cancelCtx.SetCancelChannel(common.ResetSpotTrade, stopC)
		<-doneC
		time.Sleep(100 * time.Millisecond)
	}
}

func (b *BinaMarketWs) registerFutureAggTrade() {

	wsHandler := func(event *futures.WsAggTradeEvent) {

		symbol := GetSysSymbol(event.Symbol)
		data := b.FutureAggTradeHandler(event)
		dataEvent := common.NewDataEvent(common.BINANCEID, BinanceMdWsAccountIdx, common.FutureID, common.AggTradeID, symbol, data)
		b.systemAgent.EnQueue(symbol, dataEvent)
	}

	for {

		ri := b.systemAgent.StrategyManagerCfg().EXWsCfg[common.BINANCEID].RegisterInfo.RegisterFutureAggTrade
		if len(ri) == 0 {
			stopC := make(chan struct{})
			b.cancelCtx.SetCancelChannel(common.ResetAggTrade, stopC)
			<-stopC
			time.Sleep(100 * time.Millisecond)
			continue
		}

		common.Logger.Info(b.name + " Establishing Future aggTrade connection. Please wait")
		doneC, stopC, err := futures.WsCombinedAggTradeServe(ConvertToExchangeSymbolsSlice(ri), wsHandler, b.errHandler)
		if err != nil {
			common.Logger.Error(b.name+" Future aggTrade Websocket Error:", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		b.cancelCtx.SetCancelChannel(common.ResetAggTrade, stopC)
		<-doneC
		time.Sleep(100 * time.Millisecond)
	}
}

func (b *BinaMarketWs) ResetWs(rs []common.ResetID) {
	fmt.Println("BinaMarketWs.ResetWs  ", rs)

	b.cancelCtx.CloseCancelChannels(rs)
}
