package common

import "sync/atomic"

type TickTrigInfo struct {
	PrevTime int64
	PrevBap  float64
	PrevBbp  float64
}

type RegisterWsConfig struct {
	RegisterWs       []map[string]string `json:"register_ws"`
	TrigInterval     int64               `json:"trig_interval"`
	TickTrigInterval int64               `json:"tick_trig_interval"`
	KlineInterval    string              `json:"kline_interval"`
	Options          map[string]string   `json:"options"`
}

// StrategyConfig
type StrategyCfg struct {
	KeyMap          map[string]string                           `json:"key_map"`
	OnTimerInterval int64                                       `json:"timer_interval"`
	GatewayConfigs  map[ExchangeID]map[string]*RegisterWsConfig `json:"gateway_configs"`
}

type StrategyManagerCfg struct {
	GatewayConfigs map[ExchangeID]map[string]*RegisterWsConfig `json:"gateway_configs"`
}

type IntervalInfo struct {
	TrigInterval     int64
	TrigOnOff        atomic.Value
	TickTrigInterval int64 //是booktick
	TickTrigOnOff    atomic.Value
	KlineInterval    string
}

type WsInfo struct {
	AccountWs AccountWsInterface
	MarktetWs MarketWsInterface
}
