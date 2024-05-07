package common

type BinanceRestClientInterface interface {
	GetDustAssets() (res *ListDustResponse, err error)
	ConvertDustAssets(assets []SymbolID) (*DustTransferResponse, error)
	GetMarginDustAssets() (*[]ListMarginDustResponse, error)
	ConvertMarginDustAssets(assets []SymbolID) (*[]ListMarginDustResponse, error)
	MarginLoan(asset SymbolID, isIsolated bool, symbol SymbolID, amount float64) (*TransactionResponse, error)
	MarginRepay(asset SymbolID, isIsolated bool, symbol SymbolID, amount float64) (*TransactionResponse, error)
	MarginAllAssets() ([]*MarginAsset, error)
	MarginAllPairs() ([]*MarginAllPair, error)
	CrossMarginCollateralRatio() ([]*CrossMarginCollateralRatio, error)
	NextHourlyInterestRates(assets []SymbolID, isIsolated bool) ([]*NextHourlyInterestRate, error)
}

type Filters struct {
	LotSizeFilter
	PriceFilter
	FilterType string `json:"filterType,omitempty"`
}

// ExchangeInfo exchange info
type ExchangeInfo struct {
	Timezone   string      `json:"timezone"`
	ServerTime int64       `json:"serverTime"`
	RateLimits []RateLimit `json:"rateLimits"`
	// ExchangeFilters []interface{}    `json:"exchangeFilters"`
	Symbols []BinaSymbolInfo `json:"symbols"`
}

// RateLimit struct
type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
}

// Symbol market symbol
type BinaSymbolInfo struct {
	Symbol                     SymbolID  `json:"symbol"`
	Status                     string    `json:"status"`
	BaseAsset                  SymbolID  `json:"baseAsset"`
	BaseAssetPrecision         int       `json:"baseAssetPrecision"`
	QuoteAsset                 SymbolID  `json:"quoteAsset"`
	QuotePrecision             int       `json:"quotePrecision"`
	QuoteAssetPrecision        int       `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int32     `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int32     `json:"quoteCommissionPrecision"`
	OrderTypes                 []string  `json:"orderTypes"`
	IcebergAllowed             bool      `json:"icebergAllowed"`
	OcoAllowed                 bool      `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool      `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool      `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool      `json:"isMarginTradingAllowed"`
	Filters                    []Filters `json:"filters"`
	Permissions                []string  `json:"permissions"`
}

// LotSizeFilter define lot size filter of symbol
type LotSizeFilter struct {
	MaxQuantity string `json:"maxQty,omitempty"`
	MinQuantity string `json:"minQty,omitempty"`
	StepSize    string `json:"stepSize,omitempty"`
}

// PriceFilter define price filter of symbol
type PriceFilter struct {
	MaxPrice string `json:"maxPrice,omitempty"`
	MinPrice string `json:"minPrice,omitempty"`
	TickSize string `json:"tickSize,omitempty"`
}

// LotSizeFilter return lot size filter of symbol
func (s *BinaSymbolInfo) LotSizeFilter() *LotSizeFilter {
	for _, filter := range s.Filters {
		if filter.FilterType == "LOT_SIZE" {
			return &filter.LotSizeFilter
		}
	}
	return nil
}

// PriceFilter return price filter of symbol
func (s *BinaSymbolInfo) PriceFilter() *PriceFilter {
	for _, filter := range s.Filters {
		if filter.FilterType == "PRICE_FILTER" {
			return &filter.PriceFilter
		}
	}
	return nil
}

type NextHourlyInterestRate struct {
	Name               SymbolID `json:"asset"`
	NextHourlyInterest string   `json:"nextHourlyInterestRate"`
}

type Collateral struct {
	MinUsdValue  string `json:"minUsdValue,omitempty"`
	MaxUsdValue  string `json:"maxUsdValue,omitempty"`
	DiscountRate string `json:"discountRate"`
}

type CrossMarginCollateralRatio struct {
	Collaterals []Collateral `json:"collaterals"`
	AssetNames  []SymbolID   `json:"assetNames"`
}
type MarginAllPair struct {
	ID            int64    `json:"id"`
	Symbol        SymbolID `json:"symbol"`
	Base          string   `json:"base"`
	Quote         string   `json:"quote"`
	IsMarginTrade bool     `json:"isMarginTrade"`
	IsBuyAllowed  bool     `json:"isBuyAllowed"`
	IsSellAllowed bool     `json:"isSellAllowed"`
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
	Asset           SymbolID `json:"asset"`
	Interest        string   `json:"interest"`
	Principal       string   `json:"principal"`
	LiabilityOfBUSD string   `json:"liabilityOfBUSD"`
}

type ListDustDetail struct {
	Asset            SymbolID `json:"asset"`
	AssetFullName    string   `json:"assetFullName"`
	AmountFree       string   `json:"amountFree"`
	ToBTC            string   `json:"toBTC"`
	ToBNB            string   `json:"toBNB"`
	ToBNBOffExchange string   `json:"toBNBOffExchange"`
	Exchange         string   `json:"exchange"`
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
	Amount              string   `json:"amount"`
	FromAsset           SymbolID `json:"fromAsset"`
	OperateTime         int64    `json:"operateTime"`
	ServiceChargeAmount string   `json:"serviceChargeAmount"`
	TranID              int64    `json:"tranId"`
	TransferedAmount    string   `json:"transferedAmount"`
}
