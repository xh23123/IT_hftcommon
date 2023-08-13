package unimargin

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetBalanceService get account balance
type GetBalanceService struct {
	c *Client
}

// Do send request
func (s *GetBalanceService) Do(ctx context.Context, opts ...RequestOption) (res []*UnimarginBalance, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/balance",
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*UnimarginBalance{}, err
	}
	res = make([]*UnimarginBalance, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UnimarginBalance{}, err
	}
	return res, nil
}

/*
   "asset": "USDT",    // 资产
    "": "122607.35137903", // 钱包余额 =  全仓杠杆未锁定 + 全仓杠杆锁定 + u本位合约钱包余额 + 币本位合约钱包余额
    "": "92.27530794", // 全仓资产 = 全仓杠杆未锁定 + 全仓杠杆锁定
    "": "10.00000000", // 全仓杠杆借贷
    "": "100.00000000", // 全仓杠杆未锁定
    "": "0.72469206", // 全仓杠杆利息
    "": "3.00000000", //全仓杠杆锁定
    "": "92.27530794", // 全仓杠杆净资产
    "": "0.00000000",  // u本位合约钱包余额
    "": "23.72469206",     // u本位未实现盈亏
    "": "23.72469206",       // 币本位合约钱包余额
    "": "",    // 币本位未实现盈亏
    "": 1617939110373
*/
// UnimarginBalance define user balance of your account
type UnimarginBalance struct {
	Asset               string `json:"asset"`
	TotalWalletBalance  string `json:"totalWalletBalance"`
	CrossMarginAsset    string `json:"crossMarginAsset"`
	CrossMarginBorrowed string `json:"crossMarginBorrowed"`
	CrossMarginFree     string `json:"crossMarginFree"`
	CrossMarginInterest string `json:"crossMarginInterest"`
	CrossMarginLocked   string `json:"crossMarginLocked"`
	UmWalletBalance     string `json:"umWalletBalance"`
	UmUnrealizedPNL     string `json:"umUnrealizedPNL"`
	CmWalletBalance     string `json:"cmWalletBalance"`
	CmUnrealizedPNL     string `json:"cmUnrealizedPNL"`
	UpdateTime          int64  `json:"updateTime"`
}

// GetPositionService get account balance
type GetUmPositionService struct {
	c *Client
}

// Do send request
func (s *GetUmPositionService) Do(ctx context.Context, opts ...RequestOption) (res []*UnimarginPosition, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/um/positionRisk",
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*UnimarginPosition{}, err
	}
	res = make([]*UnimarginPosition, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UnimarginPosition{}, err
	}
	return res, nil
}

/*
   "": "0.00000", // 开仓均价
   "": "10", // 当前杠杆倍数
   "": "6679.50671178",   // 当前标记价格
   "": "20000000", // 当前杠杆倍数允许的名义价值上限
   "": "0.000", // 头寸数量，符号代表多空方向, 正数为多，负数为空
   "": "0",
   "": "BTCUSDT", // 交易对
   "": "0.00000000", // 持仓未实现盈亏
   "": "BOTH", // 持仓方向
   "": 1625474304765   // 更新时间

*/
// UnimarginPosition define user balance of your account
type UnimarginPosition struct {
	EntryPrice       string           `json:"entryPrice"`
	Leverage         string           `json:"leverage"`
	MarkPrice        string           `json:"markPrice"`
	MaxNotionalValue string           `json:"maxNotionalValue"`
	PositionAmt      string           `json:"positionAmt"`
	Notional         string           `json:"notional"`
	Symbol           string           `json:"symbol"`
	UnRealizedProfit string           `json:"unRealizedProfit"`
	PositionSide     PositionSideType `json:"positionSide"`
	UpdateTime       int64            `json:"updateTime"`
}

// GetPositionService get account balance
type GetCmPositionService struct {
	c *Client
}

// Do send request
func (s *GetCmPositionService) Do(ctx context.Context, opts ...RequestOption) (res []*UnimarginPosition, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/cm/positionRisk",
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)

	if err != nil {
		return []*UnimarginPosition{}, err
	}
	res = make([]*UnimarginPosition, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UnimarginPosition{}, err
	}
	return res, nil
}

// GetAccountService get account info
type GetAccountService struct {
	c *Client
}

// Do send request
func (s *GetAccountService) Do(ctx context.Context, opts ...RequestOption) (res *Account, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/account",
		secType:  secTypeSigned,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Account)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Account define account info
type Account struct {
	Assets      []*AccountAsset    `json:"assets"`
	CanDeposit  bool               `json:"canDeposit"`
	CanTrade    bool               `json:"canTrade"`
	CanWithdraw bool               `json:"canWithdraw"`
	FeeTier     int                `json:"feeTier"`
	Positions   []*AccountPosition `json:"positions"`
	UpdateTime  int64              `json:"updateType"`
}

// AccountAsset define account asset
type AccountAsset struct {
	Asset                  string `json:"asset"`
	WalletBalance          string `json:"walletBalance"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	MarginBalance          string `json:"marginBalance"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	AvailableBalance       string `json:"availableBalance"`
}

// AccountPosition define accoutn position
type AccountPosition struct {
	Symbol                 string           `json:"symbol"`
	PositionAmt            string           `json:"positionAmt"`
	InitialMargin          string           `json:"initialMargin"`
	MaintMargin            string           `json:"maintMargin"`
	UnrealizedProfit       string           `json:"unrealizedProfit"`
	PositionInitialMargin  string           `json:"positionInitialMargin"`
	OpenOrderInitialMargin string           `json:"openOrderInitialMargin"`
	Leverage               string           `json:"leverage"`
	Isolated               bool             `json:"isolated"`
	PositionSide           PositionSideType `json:"positionSide"`
	EntryPrice             string           `json:"entryPrice"`
	MaxQty                 string           `json:"maxQty"`
}
