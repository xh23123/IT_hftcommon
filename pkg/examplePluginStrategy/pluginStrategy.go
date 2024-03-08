package main

import (
	"encoding/json"
	"fmt"

	. "github.com/xh23123/IT_hftcommon/pkg/common"
)

type TestStrategy struct {
	systemAgent   TradeSystemAgent
	orderPair     SymbolID
	orderSequence int64
	orderExchange ExchangeID
	maxOrderNum   int64
}

func NewStrategy(systemAgent TradeSystemAgent) StrategyInterface {

	strategy := TestStrategy{
		orderPair:     "ETH-USDT",
		orderExchange: BINANCEID,
		systemAgent:   systemAgent,
		orderSequence: 0,
		maxOrderNum:   1}

	strategy.initMdConfig()

	return &strategy
}

func (s *TestStrategy) initMdConfig() {
	strategyConfig := StrategyCfg{
		MarketDataConfigs: MarketDataConfigs{
			MdConfigs: map[ExchangeID]map[SymbolID][]*MarketDataConfig{
				BINANCEID: {
					s.orderPair: {
						{
							TransactionId:  SpotID,
							MdCallBackName: "OnBookTick",
							BookTickOptions: &BookTickOptions{
								TrigInterval: 1,
							},
						},
					},
				},
				OKEXID: {
					s.orderPair: {
						{
							TransactionId:  SpotID,
							MdCallBackName: "OnBookTick",
							BookTickOptions: &BookTickOptions{
								TrigInterval: 1,
							},
						},
					},
				},
			},
		},
		OnTimerInterval: 5000,
	}
	s.systemAgent.InitMdConfig(&strategyConfig)
}

func (s *TestStrategy) InitPara() {
}

func (s *TestStrategy) InitVar() {
}

func (s *TestStrategy) OnExit() {
	fmt.Println("strategy exited")
}

func (s *TestStrategy) InitMyStrategy() (actions []*ActionEvent) {
	return actions
}

func (s *TestStrategy) OnBookTick(event *BookTickWs) (actions []*ActionEvent) {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnBookTick", string(eventStr))

	return actions
}

func (s *TestStrategy) OnDepth(event *DepthWs) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnDepth", string(eventStr))
	return nil
}

func (s *TestStrategy) OnTick(event *TickWs) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnTick", string(eventStr))
	return nil
}

func (s *TestStrategy) OnKlineWs(event *KlineWs) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnKlineWs", string(eventStr))
	return nil
}

func (s *TestStrategy) OnOrderbook(event *Orderbook) []*ActionEvent {
	return nil
}

func (s *TestStrategy) OnTradeWs(event *TradeWs) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnTradeWs", string(eventStr))
	return nil
}

func (s *TestStrategy) OnMarkPrice(event *MarkPriceWs) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnMarkPrice", string(eventStr))
	return nil
}

func (s *TestStrategy) OnOrder(event *OrderTradeUpdateInfo) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnOrder", string(eventStr))
	return nil
}

func (s *TestStrategy) OnTrade(event *OrderTradeUpdateInfo) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnTrade", string(eventStr))
	return nil
}

func (s *TestStrategy) OnError(event *ErrorMsg) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnError", string(eventStr))
	return nil
}

func (s *TestStrategy) OnDexBookTicks(event *DexBookTicks) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnDexBookTicks", string(eventStr))
	return nil
}

func (s *TestStrategy) OnDexTrades(event *DexTrades) []*ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnDexTrades", string(eventStr))
	return nil
}

func (s *TestStrategy) OnTimer() (actions []*ActionEvent) {

	return actions
}
