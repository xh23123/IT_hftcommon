package common

type TickTrigInfo struct {
	PrevTime int64
	PrevBap  float64
	PrevBbp  float64
}

type IntervalInfo struct {
	TickTrigInterval int64  `json:"trig_interval"` //是booktick
	KlineInterval    string `json:"kline_interval"`
}

type RegisterWsConfig struct {
	IntervalInfo
	RegisterWs map[DataID]string `json:"register_ws"`
	Options    map[string]string `json:"options"`
}

type MarketDataConfigs struct {
	MdConfigs map[ExchangeID]map[TransactionID]map[SymbolID]*RegisterWsConfig `json:"gateway_configs"`
}

// StrategyConfig
type StrategyCfg struct {
	MarketDataConfigs
	KeyMap          map[string]string `json:"key_map"`
	OnTimerInterval int64             `json:"timer_interval"`
}

type WsInfo struct {
	AccountWs AccountWsInterface
	MarktetWs MarketWsInterface
}
