package common

import "math/big"

type OrderTradeUpdateInfo struct {
	DataID          DataID        `json:"dataid"`
	Transaction     TransactionID `json:"transactionid"`
	Exchange        ExchangeID    `json:"exid"`
	AccountIndex    AccountIdx    `json:"accountidx"`
	Status          OrderStatusID `json:"status"`
	Symbol          string        `json:"symbol"`
	Id              string        `json:"id"`
	Cid             string        `json:"cid"`
	Side            SideID        `json:"side"`
	PositionSide    PositionID    `json:"position_side"`
	Type            OrderTypeID   `json:"type"`
	Size            float64       `json:"size"`
	FilledSize      float64       `json:"filled_size"`
	Price           float64       `json:"price"`
	AvgPrice        float64       `json:"avg_price"`
	LastFilledPrice float64       `json:"last_filled_size"`
	FeeAsset        string        `json:"fee_asset"`
	FeeCost         float64       `json:"fee_cost"`
	Timestamp       int64         `json:"timestamp"`
}

type Balance struct {
	Balance          float64 `json:"wb"`
	MarginBalance    float64 `json:"mb"`
	AvailableBalance float64 `json:"ab"`
	DexBalance       big.Int `json:"dex_balance"`
	DexDecimal       int64   `json:"dex_decimal"`
}

type UserAsset struct {
	Asset    string  `json:"asset"`
	Borrowed float64 `json:"borrowed"`
	Free     float64 `json:"free"`
	Interest float64 `json:"interest"`
	Locked   float64 `json:"locked"`
	NetAsset float64 `json:"netAsset"`
}

type MarginBalance struct {
	MarginLevel float64     `json:"marginLevel"`
	UserAssets  []UserAsset `json:"userAssets"`
}
type SpotBalance map[string]*Balance

type SidePosition struct {
	Amount        float64 `json:"pa"`
	EntryPrice    float64 `json:"ep"`
	UnrealizedPnL float64 `json:"up"`
}
type FuturePosition struct {
	LONG  *SidePosition `json:"LONG"`
	SHORT *SidePosition `json:"SHORT"`
	BOTH  *SidePosition `json:"BOTH"`
}
type WsFutureBalance map[string]*Balance
type WsFuturePosition map[string]*FuturePosition

type PremiumIndexInfo struct {
	Symbol          string `json:"symbol"`
	MarkPrice       string `json:"markPrice"`
	LastFundingRate string `json:"lastFundingRate"`
	NextFundingTime int64  `json:"nextFundingTime"`
	Time            int64  `json:"time"`
}

type Kline struct {
	OpenTime         int64  `json:"openTime"` // TODO seconds?
	Open             string `json:"open"`
	High             string `json:"high"`
	Low              string `json:"low"`
	Close            string `json:"close"`
	Volume           string `json:"volume"`
	QuoteAssetVolume string `json:"quoteAssetVolume"`
}
