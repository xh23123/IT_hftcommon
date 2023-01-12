package common

type BaseConfig struct {
	OnTimerInterval int64 `json:"timer_interval"`
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
