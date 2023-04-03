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
