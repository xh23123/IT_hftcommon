package common

import "sync/atomic"

type BaseConfig struct {
	OnTimerInterval int64 `json:"timer_interval"`
}

type TickTrigInfo struct {
	PrevTime int64
	PrevBap  float64
	PrevBbp  float64
}

type RegisterWsConfig struct {
	RegisterWs       []string `json:"register_ws"`
	TrigInterval     int64    `json:"trig_interval"`
	TickTrigInterval int64    `json:"tick_trig_interval"`
	KlineInterval    string   `json:"kline_interval"`
}

// StrategyConfig
type StrategyCfg struct {
	KeyMap         map[string]string
	BaseConfig     BaseConfig
	RegisterConfig map[ExchangeID]map[string]*RegisterWsConfig
}

type StrategyManagerCfg struct {
	OnTimerInterval int64
	EXWsCfg         map[ExchangeID]*StrategyWsCfg
}

type StrategyWsCfg struct {
	IntervalMap  map[string]*IntervalInfo
	RegisterInfo *RegisterInfo
}

type IntervalInfo struct {
	TrigInterval     int64
	TrigOnOff        atomic.Value
	TickTrigInterval int64 //是booktick
	TickTrigOnOff    atomic.Value
	KlineInterval    string
}

type RegisterInfo struct {
	RegisterSpotDepth          map[string]string
	RegisterSpotBookTick       []string
	RegisterSpotTick           []string
	RegisterSpotKline          []string
	RegisterSpotTrade          []string
	RegisterFutureDepth        map[string]string
	RegisterFutureBookTick     []string
	RegisterFutureTick         []string
	RegisterFutureKline        []string
	RegisterCoinFutureBookTick []string
	RegisterMarkPrice          []string
	RegisterSpotKlineWs        map[string]string
	RegisterFutureKlineWs      map[string]string
	RegisterFutureAggTrade     []string
}

type WsInfo struct {
	AccountWs AccountWsInterface
	MarktetWs MarketWsInterface
}
