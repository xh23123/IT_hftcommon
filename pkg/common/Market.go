package common

// 行情相关数据
type BookTickWs struct {
	Exchange          ExchangeID `json:"exid"`
	DataID            DataID     `json:"dataid"`
	Symbol            string     `json:"symbol"`
	UpdateID          int64      `json:"updateid"`
	BestBidPrice      float64    `json:"bbp"`
	BestBidSize       float64    `json:"bbs"`
	BestAskPrice      float64    `json:"bap"`
	BestAskSize       float64    `json:"bas"`
	ReceiveTimestamp  int64      `json:"recvtimestamp"`
	ExchangeTimestamp int64      `json:"extimestamp"`
}

type DepthWs struct {
	Exchange          ExchangeID     `json:"exid"`
	DataID            DataID         `json:"dataid"`
	Symbol            string         `json:"symbol"`
	UpdateID          int64          `json:"updateid"`
	Bids              [20][2]float64 `json:"bids"`
	Asks              [20][2]float64 `json:"asks"`
	ReceiveTimestamp  int64          `json:"recvtimestamp"`
	ExchangeTimestamp int64          `json:"extimestamp"`
}

type TradeWs struct {
	Exchange          ExchangeID `json:"exid"`
	DataID            DataID     `json:"dataid"`
	Symbol            string     `json:"symbol"`
	TradeID           int64      `json:"tradeid"`
	Price             float64    `json:"price"`
	Size              float64    `json:"size"`
	TradeCount        int64      `json:"tradecount"`
	IsMaker           bool       `json:"m"`
	ReceiveTimestamp  int64      `json:"recvtimestamp"`
	ExchangeTimestamp int64      `json:"extimestamp"`
}

type TickWs struct {
	Exchange          ExchangeID `json:"exid"`
	DataID            DataID     `json:"dataid"`
	Symbol            string     `json:"symbol"`
	Price             float64    `json:"price"`
	Size              float64    `json:"size"`
	IsBuyerMaker      bool       `json:"m"`
	ReceiveTimestamp  int64      `json:"recvtimestamp"`
	ExchangeTimestamp int64      `json:"extimestamp"`
}

type MarkPriceWs struct {
	Exchange          ExchangeID `json:"exid"`
	DataID            DataID     `json:"dataid"`
	Symbol            string     `json:"symbol"`
	MarkPrice         float64    `json:"mark_price"`
	IndexPrice        float64    `json:"index_price"`
	FundingRate       float64    `json:"funding_rate"`
	ReceiveTimestamp  int64      `json:"recvtimestamp"`
	ExchangeTimestamp int64      `json:"extimestamp"`
}

type KlineWs struct {
	Exchange          ExchangeID `json:"exid"`
	DataID            DataID     `json:"dataid"`
	Symbol            string     `json:"symbol"`
	OpenTime          int64      `json:"open_time"`
	CloseTime         int64      `json:"close_time"`
	TradeCount        int64      `json:"trade_count"`
	Open              float64    `json:"open"`
	High              float64    `json:"high"`
	Low               float64    `json:"low"`
	Close             float64    `json:"close"`
	TradeSize         float64    `json:"trade_size"`
	TradeVolume       float64    `json:"trade_value"`
	BuySize           float64    `json:"buy_size"`
	BuyVolume         float64    `json:"buy_value"`
	ReceiveTimestamp  int64      `json:"recvtimestamp"`
	ExchangeTimestamp int64      `json:"extimestamp"`
}
