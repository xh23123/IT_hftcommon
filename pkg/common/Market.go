package common

import "math/big"

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

type DexBookTick struct {
	Symbol            string             `json:"symbol"`
	ReceiveTimestamp  int64              `json:"recvtimestamp"`
	ExchangeTimestamp int64              `json:"extimestamp"`
	UniswapV2         *UniswapV2BookTick `json:"uniswapv2,omitempty"`
}

type UniswapV2BookTick struct {
	Token0   string  `json:"token0"`
	Toekn1   string  `json:"token1"`
	Decimal0 int64   `json:"decimal0"`
	Decimal1 int64   `json:"decimal1"`
	Reserve0 big.Int `json:"reserve0"`
	Reserve1 big.Int `json:"reserve1"`
}

type DexBookTicks struct {
	Exchange ExchangeID    `json:"exid"`
	DataID   DataID        `json:"dataid"`
	Ticks    []DexBookTick `json:"ticks"`
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

type Orderbook struct {
	Exchange          ExchangeID   `json:"exid"`
	DataID            DataID       `json:"dataid"`
	Symbol            string       `json:"symbol"`
	UpdateID          int64        `json:"updateid"`
	Bids              [][2]float64 `json:"bids"`
	Asks              [][2]float64 `json:"asks"`
	ReceiveTimestamp  int64        `json:"recvtimestamp"`
	ExchangeTimestamp int64        `json:"extimestamp"`
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

type DexTrade struct {
	Symbol            string          `json:"symbol"`
	ReceiveTimestamp  int64           `json:"recvtimestamp"`
	ExchangeTimestamp int64           `json:"extimestamp"`
	UniswapV2         *UniswapV2Trade `json:"uniswapv2,omitempty"`
}

type UniswapV2Trade struct {
	Token0     string  `json:"token0"`
	Toekn1     string  `json:"token1"`
	Decimal0   int64   `json:"decimal0"`
	Decimal1   int64   `json:"decimal1"`
	Amount0In  big.Int `json:"amount0_in"`
	Amount1In  big.Int `json:"amount1_in"`
	Amount0Out big.Int `json:"amount0_out"`
	Amount1Out big.Int `json:"amount1_out"`
	Removed    bool    `json:"removed"`
}

type DexTrades struct {
	Exchange ExchangeID `json:"exid"`
	DataID   DataID     `json:"dataid"`
	Trades   []DexTrade `json:"trades"`
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
