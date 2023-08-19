package common

type BinanceRestClientInterface interface {
	GetDustAssets() (res *ListDustResponse, err error)
	ConvertDustAssets(assets []string) (*DustTransferResponse, error)
	GetMarginDustAssets() (*[]ListMarginDustResponse, error)
	ConvertMarginDustAssets(assets []string) (*[]ListMarginDustResponse, error)
	MarginLoan(asset string, isIsolated bool, symbol string, amount float64) (*TransactionResponse, error)
	MarginRepay(asset string, isIsolated bool, symbol string, amount float64) (*TransactionResponse, error)
	MarginAllAssets() ([]*MarginAsset, error)
	MarginAllPairs() ([]*MarginAllPair, error)
	CrossMarginCollateralRatio() ([]*CrossMarginCollateralRatio, error)
	NextHourlyInterestRates(assets []string, isIsolated bool) ([]*NextHourlyInterestRate, error)
	GetExchangeInfos(transactionId TransactionID) (*ExchangeInfo, error)
}

// ExchangeInfo exchange info
type ExchangeInfo struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int64         `json:"serverTime"`
	RateLimits      []RateLimit   `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []Symbol      `json:"symbols"`
}

// RateLimit struct
type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
}

// Symbol market symbol
type Symbol struct {
	Symbol                     string                   `json:"symbol"`
	Status                     string                   `json:"status"`
	BaseAsset                  string                   `json:"baseAsset"`
	BaseAssetPrecision         int                      `json:"baseAssetPrecision"`
	QuoteAsset                 string                   `json:"quoteAsset"`
	QuotePrecision             int                      `json:"quotePrecision"`
	QuoteAssetPrecision        int                      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int32                    `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int32                    `json:"quoteCommissionPrecision"`
	OrderTypes                 []string                 `json:"orderTypes"`
	IcebergAllowed             bool                     `json:"icebergAllowed"`
	OcoAllowed                 bool                     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool                     `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool                     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool                     `json:"isMarginTradingAllowed"`
	Filters                    []map[string]interface{} `json:"filters"`
	Permissions                []string                 `json:"permissions"`
}

// LotSizeFilter define lot size filter of symbol
type LotSizeFilter struct {
	MaxQuantity string `json:"maxQty"`
	MinQuantity string `json:"minQty"`
	StepSize    string `json:"stepSize"`
}

// PriceFilter define price filter of symbol
type PriceFilter struct {
	MaxPrice string `json:"maxPrice"`
	MinPrice string `json:"minPrice"`
	TickSize string `json:"tickSize"`
}

// LotSizeFilter return lot size filter of symbol
func (s *Symbol) LotSizeFilter() *LotSizeFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string("LOT_SIZE") {
			f := &LotSizeFilter{}
			if i, ok := filter["maxQty"]; ok {
				f.MaxQuantity = i.(string)
			}
			if i, ok := filter["minQty"]; ok {
				f.MinQuantity = i.(string)
			}
			if i, ok := filter["stepSize"]; ok {
				f.StepSize = i.(string)
			}
			return f
		}
	}
	return nil
}

// PriceFilter return price filter of symbol
func (s *Symbol) PriceFilter() *PriceFilter {
	for _, filter := range s.Filters {
		if filter["filterType"].(string) == string("PRICE_FILTER") {
			f := &PriceFilter{}
			if i, ok := filter["maxPrice"]; ok {
				f.MaxPrice = i.(string)
			}
			if i, ok := filter["minPrice"]; ok {
				f.MinPrice = i.(string)
			}
			if i, ok := filter["tickSize"]; ok {
				f.TickSize = i.(string)
			}
			return f
		}
	}
	return nil
}

type NextHourlyInterestRate struct {
	Name               string `json:"asset"`
	NextHourlyInterest string `json:"nextHourlyInterestRate"`
}

type Collateral struct {
	MinUsdValue  string `json:"minUsdValue,omitempty"`
	MaxUsdValue  string `json:"maxUsdValue,omitempty"`
	DiscountRate string `json:"discountRate"`
}

type CrossMarginCollateralRatio struct {
	Collaterals []Collateral `json:"collaterals"`
	AssetNames  []string     `json:"assetNames"`
}
type MarginAllPair struct {
	ID            int64  `json:"id"`
	Symbol        string `json:"symbol"`
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	IsMarginTrade bool   `json:"isMarginTrade"`
	IsBuyAllowed  bool   `json:"isBuyAllowed"`
	IsSellAllowed bool   `json:"isSellAllowed"`
}

type MarginAsset struct {
	FullName      string `json:"assetFullName"`
	Name          string `json:"assetName"`
	Borrowable    bool   `json:"isBorrowable"`
	Mortgageable  bool   `json:"isMortgageable"`
	UserMinBorrow string `json:"userMinBorrow"`
	UserMinRepay  string `json:"userMinRepay"`
}

type TransactionResponse struct {
	TranID int64 `json:"tranId"`
}

type ListMarginDustResponse struct {
	Asset           string `json:"asset"`
	Interest        string `json:"interest"`
	Principal       string `json:"principal"`
	LiabilityOfBUSD string `json:"liabilityOfBUSD"`
}

type ListDustDetail struct {
	Asset            string `json:"asset"`
	AssetFullName    string `json:"assetFullName"`
	AmountFree       string `json:"amountFree"`
	ToBTC            string `json:"toBTC"`
	ToBNB            string `json:"toBNB"`
	ToBNBOffExchange string `json:"toBNBOffExchange"`
	Exchange         string `json:"exchange"`
}

type ListDustResponse struct {
	Details            []ListDustDetail `json:"details"`
	TotalTransferBtc   string           `json:"totalTransferBtc"`
	TotalTransferBNB   string           `json:"totalTransferBNB"`
	DribbletPercentage string           `json:"dribbletPercentage"`
}

// DustTransferResponse represents the response from DustTransferService.
type DustTransferResponse struct {
	TotalServiceCharge string                `json:"totalServiceCharge"`
	TotalTransfered    string                `json:"totalTransfered"`
	TransferResult     []*DustTransferResult `json:"transferResult"`
}

// DustTransferResult represents the result of a dust transfer.
type DustTransferResult struct {
	Amount              string `json:"amount"`
	FromAsset           string `json:"fromAsset"`
	OperateTime         int64  `json:"operateTime"`
	ServiceChargeAmount string `json:"serviceChargeAmount"`
	TranID              int64  `json:"tranId"`
	TransferedAmount    string `json:"transferedAmount"`
}
