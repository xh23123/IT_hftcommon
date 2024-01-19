package common

type BookTickOptions struct {
	TrigInterval int64 `json:"trig_interval"`
}

type KlineOptions struct {
	KlineInterval string `json:"kline_interval"`
}

type OrderBookOptions struct {
	TopLevels int `json:"top_levels"`
}

type MarketDataConfig struct {
	TransactionId  TransactionID `json:"transaction_id"`
	MdCallBackName string        `json:"md_callback_name"`

	BookTickOptions  *BookTickOptions  `json:"book_tick_options"`
	KlineOptions     *KlineOptions     `json:"kline_options"`
	OrderBookOptions *OrderBookOptions `json:"order_book_options"`
}

type MarketDataConfigs struct {
	MdConfigs map[ExchangeID]map[SymbolID][]*MarketDataConfig `json:"md_configs"`
}

// StrategyConfig
type StrategyCfg struct {
	MarketDataConfigs
	KeyMap          map[SymbolID]string `json:"key_map"`
	OnTimerInterval int64               `json:"timer_interval"`
}

type WsInfo struct {
	AccountWs AccountWsInterface
	MarktetWs MarketWsInterface
}
