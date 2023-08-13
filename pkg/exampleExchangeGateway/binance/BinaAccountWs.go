package bina

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xh23123/IT_hftcommon/pkg/common"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2/delivery"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2/futures"

	"go.uber.org/zap"
)

var _ common.AccountWsInterface = (*BinaAccountWs)(nil)

type BinaAccountWs struct {
	systemAgent     common.TradeSystemAgent
	accountIndex    common.AccountIdx
	spotCli         *binance.Client
	futureCli       *futures.Client
	futureListenKey string
	spotListenKey   string

	coinFutureCli *delivery.Client
	name          string
	invertPairMap map[string]map[string]bool

	marginListenKey string

	coinFutureListenKey string
	cancelCtx           map[common.ResetID]chan struct{}
}

func NewAccountWs(systemAgent common.TradeSystemAgent, info AccountInfo) {

	baw := &BinaAccountWs{
		systemAgent:  systemAgent,
		accountIndex: info.Index,
		name:         fmt.Sprint("binanceAccount:", info.Index)}
	apiKey := systemAgent.Config().Section(info.Name).Key("api_key").String()
	secretKey := systemAgent.Config().Section(info.Name).Key("secret_key").String()

	baw.spotCli = binance.NewClient(apiKey, secretKey)
	baw.futureCli = futures.NewClient(apiKey, secretKey)
	baw.invertPairMap = map[string]map[string]bool{}
	baw.cancelCtx = map[common.ResetID]chan struct{}{
		common.ResetSpotAccount:       nil,
		common.ResetFutureAccount:     nil,
		common.ResetCoinFutureAccount: nil,
	}

	systemAgent.RegisterAccountWs(common.BINANCEID, info.Index, baw)
	baw.registerAcccountWs()
	go baw.onTimer()

}

func (b *BinaAccountWs) registerAcccountWs() {
	go b.registerSpotAcccountWs()
	go b.registerFutureAcccountWs()
}

func (b *BinaAccountWs) spotWsHandler(event *binance.WsUserDataEvent) {

	switch event.Event {
	case binance.UserDataEventTypeOutboundAccountPosition:
		// 重组asset所代表的symbol
		symbol := ""
		asset := []string{}
		for _, a := range event.AccountUpdate.WsAccountUpdates {
			if a.Asset == "BNB" && len(event.AccountUpdate.WsAccountUpdates) == 3 {
				continue
			}
			asset = append(asset, a.Asset)
		}
		if len(asset) == 0 {
			//这种消息直接丢弃
			common.Logger.Warn(b.name+" Wired Data", zap.Any("event", event))
			return
		}
		for k := range b.invertPairMap[asset[0]] {
			if _, ok := b.invertPairMap[asset[1]][k]; ok {
				symbol = GetSysSymbol(k)
			}
		}
		// 如果这个symbol在gateway中顺序执行
		// 否则直接只执行
		data := b.accountUpdateDo(event)
		b.systemAgent.EnQueue(symbol, common.NewDataEvent(common.BINANCEID, b.accountIndex, common.SpotID, common.AccountUpdateID, symbol, data))
	case binance.UserDataEventTypeExecutionReport:
		symbol := GetSysSymbol(event.OrderUpdate.Symbol)
		//trade类型
		switch event.OrderUpdate.Status {
		case "FILLED", "PARTIALLY_FILLED":

			tradeUpdate := b.tradeUpdateDo(symbol, common.SpotID, event)
			b.systemAgent.FeedbackTimestamp(common.BINANCEID,
				b.accountIndex,
				common.TradeUpdateID,
				symbol,
				common.UNKNOWN_ACTION,
				common.OnTradeToGateway,
				tradeUpdate.Cid,
				tradeUpdate.Id,
				tradeUpdate.Timestamp,
			)
			b.systemAgent.EnQueue(symbol, common.NewDataEvent(common.BINANCEID, b.accountIndex, common.SpotID, common.TradeUpdateID, symbol, *tradeUpdate))
		default:

			//order update
			orderUpdate := b.orderUpdateDo(symbol, common.SpotID, event)
			b.systemAgent.FeedbackTimestamp(common.BINANCEID,
				b.accountIndex,
				common.OrderUpdateID,
				symbol,
				common.UNKNOWN_ACTION,
				common.OnOrderToGateway,
				orderUpdate.Cid,
				orderUpdate.Id,
				orderUpdate.Timestamp,
			)

			b.systemAgent.EnQueue(symbol, common.NewDataEvent(common.BINANCEID, b.accountIndex, common.SpotID, common.OrderUpdateID, symbol, *orderUpdate))
		}
	case binance.UserDataEventTypeBalanceUpdate:
	default:
		message, _ := json.Marshal(event)
		common.Logger.Warn(b.name + " Spot AccountWs Unknown event, message: " + string(message))
	}
	common.Logger.Info(b.name+" spotWsHandler event: ", zap.Any("event", event))
}

func (b *BinaAccountWs) futureWsHandler(event *futures.WsUserDataEvent) {
	switch event.Event {
	case "ACCOUNT_UPDATE":
		data := b.futureAccountUpdateDo(event)
		if lens := len(event.AccountUpdate.Positions); lens > 0 {
			symbol, err := TryGetSysSymbol(event.AccountUpdate.Positions[0].Symbol)

			if err != nil {
				message, _ := json.Marshal(event)
				common.Logger.Error(b.name+" futureWsHandler unknown ACCOUNT_UPDATE symbol ",
					zap.String("symbol", event.OrderTradeUpdate.Symbol),
					zap.String("message", string(message)))

				panic("bina GetSysSymbol cant find " + event.OrderTradeUpdate.Symbol)
			}
			b.systemAgent.EnQueue(symbol, common.NewDataEvent(common.BINANCEID, b.accountIndex, common.FutureID, common.AccountUpdateID, symbol, data))
		} else {
			b.systemAgent.WsUpdateFutureBalancePosition(common.BINANCEID, b.accountIndex, data)
		}
	case "ORDER_TRADE_UPDATE":
		symbol, err := TryGetSysSymbol(event.OrderTradeUpdate.Symbol)
		if err != nil {
			message, _ := json.Marshal(event)
			common.Logger.Error(b.name+" futureWsHandler unknown ORDER_TRADE_UPDATE symbol ",
				zap.String("symbol", event.OrderTradeUpdate.Symbol),
				zap.String("message", string(message)))

			panic("bina GetSysSymbol cant find " + event.OrderTradeUpdate.Symbol)
		}

		switch event.OrderTradeUpdate.Status {
		case "FILLED", "PARTIALLY_FILLED":
			tradeUpdate := b.futureTradeUpdateDo(symbol, event)
			b.systemAgent.FeedbackTimestamp(common.BINANCEID,
				b.accountIndex,
				common.TradeUpdateID,
				symbol,
				common.UNKNOWN_ACTION,
				common.OnTradeToGateway,
				tradeUpdate.Cid,
				tradeUpdate.Id,
				tradeUpdate.Timestamp,
			)
			b.systemAgent.EnQueue(symbol, common.NewDataEvent(common.BINANCEID, b.accountIndex, common.FutureID, common.TradeUpdateID, symbol, *tradeUpdate))
		default:
			//order update
			orderUpdate := b.futureOrderUpdateDo(symbol, event)
			b.systemAgent.FeedbackTimestamp(common.BINANCEID,
				b.accountIndex,
				common.OrderUpdateID,
				symbol,
				common.UNKNOWN_ACTION,
				common.OnOrderToGateway,
				orderUpdate.Cid,
				orderUpdate.Id,
				orderUpdate.Timestamp,
			)

			b.systemAgent.EnQueue(symbol, common.NewDataEvent(common.BINANCEID, b.accountIndex, common.FutureID, common.OrderUpdateID, symbol, *orderUpdate))
		}

	default:
		message, _ := json.Marshal(event)
		common.Logger.Warn(b.name + " Future AccountWs Unknown event,message: " + string(message))
	}

	common.Logger.Info(b.name+" futureWsHandler event: ", zap.Any("event", event))
}

func (b *BinaAccountWs) registerSpotAcccountWs() {
	errHandler := func(err error) {
		common.Logger.Error(b.name+" Spot Account Websocket Error:", zap.String("error", err.Error()))
	}
	for {
		common.Logger.Info(b.name + " Establishing spot account connection. Please wait")
		// 获取listenkey
		err := b.resetSpotListenkey()
		if err != nil {
			common.Logger.Error(b.name+" Spot Account Websocket resetSpotListenkey Error:", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		// 建立用户websocket连接
		doneC, stopC, err := binance.WsUserDataServe(b.spotListenKey, b.spotWsHandler, errHandler)
		if err != nil {
			common.Logger.Error(b.name+" Spot Account Websocket WsUserDataServe Error:", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		b.cancelCtx[common.ResetSpotAccount] = stopC
		<-doneC
	}
}

func (b *BinaAccountWs) registerFutureAcccountWs() {

	errHandler := func(err error) {
		common.Logger.Error(b.name+" Future Account Websocket Error:", zap.String("error", err.Error()))
	}
	for {
		// 获取listenkey
		common.Logger.Info(b.name + " Establishing future account connection. Please wait")
		err := b.resetFutureListenkey()
		if err != nil {
			common.Logger.Error(b.name+" Future Account Websocket resetFutureListenkey Error:", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		// 建立用户websocket连接
		doneC, stopC, err := futures.WsUserDataServe(b.futureListenKey, b.futureWsHandler, errHandler)
		if err != nil {
			common.Logger.Error(b.name+" Future Account Websocket WsUserDataServe Error:", zap.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		b.cancelCtx[common.ResetFutureAccount] = stopC
		<-doneC
	}
}

func (b *BinaAccountWs) ResetWs(rs []common.ResetID) {

	common.Logger.Info(b.name + " ResetWs:" + string(fmt.Sprint("BinaAccountWs.ResetWs  ", rs)))
	for _, r := range rs {
		b.resetAccount(r)
	}
}

func (b *BinaAccountWs) resetAccount(resetId common.ResetID) {

	if b.cancelCtx[resetId] != nil {
		close(b.cancelCtx[resetId])
		b.cancelCtx[resetId] = nil
	}
}

func (b *BinaAccountWs) keepaliveSpotkey() {
	err := b.spotCli.NewKeepaliveUserStreamService().ListenKey(b.spotListenKey).Do(context.Background())
	common.Logger.Info(b.name + " Keepalive Spot ListenKey via onTimer")
	if err != nil {
		common.Logger.Error(b.name+" Keepalive Spot ListenKey via onTimer Error", zap.String("error", err.Error()))
	}
}

func (b *BinaAccountWs) keepaliveFuturekey() {
	err := b.futureCli.NewKeepaliveUserStreamService().ListenKey(b.futureListenKey).Do(context.Background())
	common.Logger.Info(b.name + " Keepalive Future ListenKey via onTimer")
	if err != nil {
		common.Logger.Error(b.name+" Keepalive Future ListenKey via onTimer Error", zap.String("error", err.Error()))
	}
}

func (b *BinaAccountWs) onTimer() {
	// 每15分钟续期listenkey
	ticker1 := time.NewTicker(15 * time.Minute)
	// 每4小时重新构建连接
	ticker2 := time.NewTicker(4 * time.Hour)
	for {
		select {
		case <-ticker1.C:
			b.keepaliveSpotkey()
			b.keepaliveFuturekey()

		case <-ticker2.C:
			common.Logger.Info(b.name + " Reset UserWs via onTimer")

			b.resetAccount(common.ResetSpotAccount)
			b.resetAccount(common.ResetFutureAccount)

		}
	}
}

func (a *BinaAccountWs) resetSpotListenkey() error {
	Listenkey, err := a.spotCli.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		return err
	}
	a.spotListenKey = Listenkey
	return nil
}

func (a *BinaAccountWs) resetFutureListenkey() error {
	Listenkey, err := a.futureCli.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		return err
	}
	a.futureListenKey = Listenkey
	return nil
}

func (a *BinaAccountWs) accountUpdateDo(event *binance.WsUserDataEvent) common.SpotBalance {
	data := common.SpotBalance{}
	for _, r := range event.AccountUpdate.WsAccountUpdates {
		f := common.Str2Float(r.Free)
		l := common.Str2Float(r.Locked)
		data[r.Asset] = &common.Balance{Balance: f + l}
	}
	return data
}

func (a *BinaAccountWs) futureAccountUpdateDo(event *futures.WsUserDataEvent) common.WsFutureBalancePosition {
	data := common.WsFutureBalancePosition{FutureBalances: make(map[string]*common.Balance),
		FuturePositions: make(map[string]*common.FuturePosition)}
	for _, r := range event.AccountUpdate.Balances {
		data.FutureBalances[r.Asset] = &common.Balance{
			Balance: common.Str2Float(r.Balance),
		}
	}
	for _, r := range event.AccountUpdate.Positions {
		symbol := GetSysSymbol(r.Symbol)
		if _, ok := data.FuturePositions[symbol]; !ok {
			data.FuturePositions[symbol] = &common.FuturePosition{}
		}
		switch r.Side {
		case futures.PositionSideTypeLong:
			data.FuturePositions[symbol].LONG = &common.SidePosition{
				Amount:        common.Str2Float(r.Amount),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedPnL),
			}
		case futures.PositionSideTypeShort:
			data.FuturePositions[symbol].SHORT = &common.SidePosition{
				Amount:        common.Str2Float(r.Amount),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedPnL),
			}
		case futures.PositionSideTypeBoth:
			data.FuturePositions[symbol].BOTH = &common.SidePosition{
				Amount:        common.Str2Float(r.Amount),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedPnL),
			}
		}
	}
	return data
}

func (a *BinaAccountWs) tradeUpdateDo(symbol string, transactionId common.TransactionID, event *binance.WsUserDataEvent) *common.OrderTradeUpdateInfo {
	data := common.OrderTradeUpdateInfo{
		DataID:       common.TradeUpdateID,
		Transaction:  transactionId,
		Exchange:     common.BINANCEID,
		AccountIndex: a.accountIndex,
		Status:       common.OrderStatusID(event.OrderUpdate.Status),
		Symbol:       symbol,
		Id:           common.Int2Str(event.OrderUpdate.Id),
		Cid:          event.OrderUpdate.ClientOrderId,
		Side:         common.SideID(event.OrderUpdate.Side),
		Type:         common.OrderTypeID(event.OrderUpdate.Type),
		Size:         common.Str2Float(event.OrderUpdate.Volume),
		FilledSize:   common.Str2Float(event.OrderUpdate.FilledVolume),
		Price:        common.Str2Float(event.OrderUpdate.Price),
		FeeAsset:     event.OrderUpdate.FeeAsset,
		FeeCost:      common.Str2Float(event.OrderUpdate.FeeCost),
		Timestamp:    time.Now().UnixNano(),
	}
	return &data
}

func (a *BinaAccountWs) orderUpdateDo(symbol string, transactionId common.TransactionID, event *binance.WsUserDataEvent) *common.OrderTradeUpdateInfo {
	cid := event.OrderUpdate.ClientOrderId
	if event.OrderUpdate.Status == "CANCELED" {
		cid = event.OrderUpdate.OrigCustomOrderId
	}
	data := common.OrderTradeUpdateInfo{
		DataID:       common.OrderUpdateID,
		Transaction:  transactionId,
		Exchange:     common.BINANCEID,
		AccountIndex: a.accountIndex,
		Status:       common.OrderStatusID(event.OrderUpdate.Status),
		Symbol:       symbol,
		Id:           common.Int2Str(event.OrderUpdate.Id),
		Cid:          cid,
		Side:         common.SideID(event.OrderUpdate.Side),
		Type:         common.OrderTypeID(event.OrderUpdate.Type),
		Size:         common.Str2Float(event.OrderUpdate.Volume),
		FilledSize:   common.Str2Float(event.OrderUpdate.FilledVolume),
		Price:        common.Str2Float(event.OrderUpdate.Price),
		FeeAsset:     event.OrderUpdate.FeeAsset,
		FeeCost:      common.Str2Float(event.OrderUpdate.FeeCost),
		Timestamp:    time.Now().UnixNano(),
	}
	return &data
}

func (a *BinaAccountWs) futureTradeUpdateDo(symbol string, event *futures.WsUserDataEvent) *common.OrderTradeUpdateInfo {
	data := common.OrderTradeUpdateInfo{
		DataID:          common.TradeUpdateID,
		Transaction:     common.FutureID,
		Exchange:        common.BINANCEID,
		AccountIndex:    a.accountIndex,
		Status:          common.OrderStatusID(event.OrderTradeUpdate.Status),
		Symbol:          symbol,
		Id:              common.Int2Str(event.OrderTradeUpdate.ID),
		Cid:             event.OrderTradeUpdate.ClientOrderID,
		Side:            common.SideID(event.OrderTradeUpdate.Side),
		PositionSide:    common.PositionID(event.OrderTradeUpdate.PositionSide),
		Type:            common.OrderTypeID(event.OrderTradeUpdate.Type),
		Size:            common.Str2Float(event.OrderTradeUpdate.OriginalQty),
		FilledSize:      common.Str2Float(event.OrderTradeUpdate.AccumulatedFilledQty),
		Price:           common.Str2Float(event.OrderTradeUpdate.OriginalPrice),
		AvgPrice:        common.Str2Float(event.OrderTradeUpdate.AveragePrice),
		LastFilledPrice: common.Str2Float(event.OrderTradeUpdate.LastFilledQty),
		Timestamp:       time.Now().UnixNano(),
	}
	return &data
}

func (a *BinaAccountWs) futureOrderUpdateDo(symbol string, event *futures.WsUserDataEvent) *common.OrderTradeUpdateInfo {
	data := common.OrderTradeUpdateInfo{
		DataID:          common.OrderUpdateID,
		Transaction:     common.FutureID,
		Exchange:        common.BINANCEID,
		AccountIndex:    a.accountIndex,
		Status:          common.OrderStatusID(event.OrderTradeUpdate.Status),
		Symbol:          symbol,
		Id:              common.Int2Str(event.OrderTradeUpdate.ID),
		Cid:             event.OrderTradeUpdate.ClientOrderID,
		Side:            common.SideID(event.OrderTradeUpdate.Side),
		PositionSide:    common.PositionID(event.OrderTradeUpdate.PositionSide),
		Type:            common.OrderTypeID(event.OrderTradeUpdate.Type),
		Size:            common.Str2Float(event.OrderTradeUpdate.OriginalQty),
		FilledSize:      common.Str2Float(event.OrderTradeUpdate.AccumulatedFilledQty),
		Price:           common.Str2Float(event.OrderTradeUpdate.OriginalPrice),
		AvgPrice:        common.Str2Float(event.OrderTradeUpdate.AveragePrice),
		LastFilledPrice: common.Str2Float(event.OrderTradeUpdate.LastFilledQty),
		Timestamp:       time.Now().UnixNano(),
	}
	return &data
}
