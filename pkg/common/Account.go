package common

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
	AvailableBalance float64 `json:"ab"`
}
type WsSpotBalance map[string]*Balance

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
