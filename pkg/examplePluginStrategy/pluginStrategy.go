package main

import (
	"encoding/json"
	"fmt"

	. "github.com/xh23123/IT_hftcommon/pkg/common"
)

type TestStrategy struct {
	systemAgent   TradeSystemAgent
	orderPair     string
	orderSequence int64
	orderExchange ExchangeID
	maxOrderNum   int64
	lastOrderCid  []string
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
	strategyConfig := StrategyCfg{}
	//基础配置
	strategyConfig.KeyMap = map[string]string{}
	strategyConfig.BaseConfig.OnTimerInterval = 5000
	strategyConfig.RegisterConfig = map[ExchangeID]map[string]*RegisterWsConfig{}
	strategyConfig.RegisterConfig[BINANCEID] = map[string]*RegisterWsConfig{}
	strategyConfig.RegisterConfig[OKEXID] = map[string]*RegisterWsConfig{}

	//策略配置
	//Binance
	strategyConfig.RegisterConfig[BINANCEID][s.orderPair] = getBinanceMdConfig()
	strategyConfig.RegisterConfig[OKEXID][s.orderPair] = getOkexMdConfig()
	s.systemAgent.InitMdConfig(&strategyConfig)
}

func getBinanceMdConfig() *RegisterWsConfig {
	//策略配置
	registerWsConfig := RegisterWsConfig{}
	registerWsConfig.RegisterWs = append(registerWsConfig.RegisterWs, "OnBookTick")
	//registerWsConfig.RegisterWs = append(registerWsConfig.RegisterWs, "OnFutureBookTick")
	//registerWsConfig.RegisterWs = append(registerWsConfig.RegisterWs, "OnFutureDepth")
	//registerWsConfig.RegisterWs = append(registerWsConfig.RegisterWs, "OnFutureAggTrade")

	registerWsConfig.TickTrigInterval = 1

	return &registerWsConfig
}

func getOkexMdConfig() *RegisterWsConfig {
	//策略配置
	registerWsConfig := RegisterWsConfig{}
	registerWsConfig.RegisterWs = append(registerWsConfig.RegisterWs, "OnBookTick")
	//registerWsConfig.RegisterWs = append(registerWsConfig.RegisterWs, "OnFutureBookTick")

	registerWsConfig.TickTrigInterval = 1

	return &registerWsConfig
}

func (s *TestStrategy) InitPara() {
}

func (s *TestStrategy) InitVar() {
}

func (s *TestStrategy) OnExit() {
	fmt.Println("strategy exited")
}

func (s *TestStrategy) InitMyStrategy() (actions []ActionEvent) {
	return actions
}

func (s *TestStrategy) OnBookTick(event BookTickWs) (actions []ActionEvent) {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnBookTick", string(eventStr))

	return actions
}
func (s *TestStrategy) genCid(triggerDataExid ExchangeID, accountIndex AccountIdx, triggerDataId DataID) string {
	orderClientID := s.systemAgent.GenOrderClientId(triggerDataExid, accountIndex, triggerDataId, s.orderSequence)
	s.lastOrderCid = append(s.lastOrderCid, orderClientID)
	s.orderSequence += 1
	return s.lastOrderCid[len(s.lastOrderCid)-1]
}

func (s *TestStrategy) OnFutureBookTick(event BookTickWs) (actions []ActionEvent) {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnFutureBookTick", string(eventStr))

	return actions
}

func (s *TestStrategy) OnDepth(event DepthWs) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnDepth", string(eventStr))
	return nil
}

func (s *TestStrategy) OnFutureDepth(event DepthWs) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnFutureDepth", string(eventStr))
	return nil
}

func (s *TestStrategy) OnTick(event TickWs) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnTick", string(eventStr))
	return nil
}

func (s *TestStrategy) OnFutureTick(event TickWs) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnFutureTick", string(eventStr))
	return nil
}

func (s *TestStrategy) OnKlineWs(event KlineWs) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnKlineWs", string(eventStr))
	return nil
}

func (s *TestStrategy) OnFutureKlineWs(event KlineWs) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnFutureKlineWs", string(eventStr))
	return nil
}

func (s *TestStrategy) OnFutureAggTrade(event AggTradeWs) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnFutureAggTrade", string(eventStr))
	return nil
}

func (s *TestStrategy) OnMarkPrice(event MarkPriceWs) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnMarkPrice", string(eventStr))
	return nil
}

func (s *TestStrategy) OnOrder(event OrderTradeUpdateInfo) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnOrder", string(eventStr))
	return nil
}

func (s *TestStrategy) OnTrade(event OrderTradeUpdateInfo) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnTrade", string(eventStr))
	return nil
}

func (s *TestStrategy) OnError(event ErrorMsg) []ActionEvent {
	eventStr, _ := json.Marshal(event)
	fmt.Println("OnError", string(eventStr))
	return nil
}

func (s *TestStrategy) OnTimer() (actions []ActionEvent) {

	return actions
}
