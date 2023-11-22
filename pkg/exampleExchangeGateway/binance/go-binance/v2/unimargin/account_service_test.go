package unimargin

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type accountServiceTestSuite struct {
	baseTestSuite
}

func TestAccountService(t *testing.T) {
	suite.Run(t, new(accountServiceTestSuite))
}

func (s *accountServiceTestSuite) TestetAccount() {
	data := []byte(`{
		"assets": [
			{
				"asset": "BTC",
				"walletBalance": "0.00241969",
				"unrealizedProfit": "0.00000000",
				"marginBalance": "0.00241969",
				"maintMargin": "0.00000000",
				"initialMargin": "0.00000000",
				"positionInitialMargin": "0.00000000",
				"openOrderInitialMargin": "0.00000000",
				"maxWithdrawAmount": "0.00241969",
				"crossWalletBalance": "0.00241969",
				"crossUnPnl": "0.00000000",
				"availableBalance": "0.00241969"
			}
		],
		"positions": [
			{
				"symbol": "BTCUSD_201225",
				"positionAmt": "0",
				"initialMargin": "0",
				"maintMargin": "0",
				"unrealizedProfit": "0.00000000",
				"positionInitialMargin": "0",
				"openOrderInitialMargin": "0",
				"leverage": "125",
				"isolated": false,
				"positionSide": "BOTH",
				"entryPrice": "0.0",
				"maxQty": "50"
			},
			{
				"symbol": "BTCUSD_201225",
				"positionAmt": "0",
				"initialMargin": "0",
				"maintMargin": "0",
				"unrealizedProfit": "0.00000000",
				"positionInitialMargin": "0",
				"openOrderInitialMargin": "0",
				"leverage": "125",
				"isolated": false,
				"positionSide": "LONG",
				"entryPrice": "0.0",
				"maxQty": "50"
			},
			{
				"symbol": "BTCUSD_201225",
				"positionAmt": "0",
				"initialMargin": "0",
				"maintMargin": "0",
				"unrealizedProfit": "0.00000000",
				"positionInitialMargin": "0",
				"openOrderInitialMargin": "0",
				"leverage": "125",
				"isolated": false,
				"positionSide": "SHORT",
				"entryPrice": "0.0",
				"maxQty": "50"
			}
		],
		"canDeposit": true,
		"canTrade": true,
		"canWithdraw": true,
		"feeTier": 2,
		"updateTime": 0
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetAccountService().Do(newContext())
	s.r().NoError(err)
	e := &Account{
		Assets: []*AccountAsset{
			{
				Asset:                  "BTC", // asset name
				WalletBalance:          "0.00241969",
				UnrealizedProfit:       "0.00000000",
				MarginBalance:          "0.00241969",
				MaintMargin:            "0.00000000",
				InitialMargin:          "0.00000000",
				PositionInitialMargin:  "0.00000000",
				OpenOrderInitialMargin: "0.00000000",
				MaxWithdrawAmount:      "0.00241969",
				CrossWalletBalance:     "0.00241969",
				CrossUnPnl:             "0.00000000",
				AvailableBalance:       "0.00241969",
			},
		},
		Positions: []*AccountPosition{
			{
				Symbol:                 "BTCUSD_201225",
				PositionAmt:            "0",
				InitialMargin:          "0",
				MaintMargin:            "0",
				UnrealizedProfit:       "0.00000000",
				PositionInitialMargin:  "0",
				OpenOrderInitialMargin: "0",
				Leverage:               "125",
				Isolated:               false,
				PositionSide:           "BOTH",
				EntryPrice:             "0.0",
				MaxQty:                 "50",
			},
			{
				Symbol:                 "BTCUSD_201225",
				PositionAmt:            "0",
				InitialMargin:          "0",
				MaintMargin:            "0",
				UnrealizedProfit:       "0.00000000",
				PositionInitialMargin:  "0",
				OpenOrderInitialMargin: "0",
				Leverage:               "125",
				Isolated:               false,
				PositionSide:           "LONG",
				EntryPrice:             "0.0",
				MaxQty:                 "50",
			},
			{
				Symbol:                 "BTCUSD_201225",
				PositionAmt:            "0",
				InitialMargin:          "0",
				MaintMargin:            "0",
				UnrealizedProfit:       "0.00000000",
				PositionInitialMargin:  "0",
				OpenOrderInitialMargin: "0",
				Leverage:               "125",
				Isolated:               false,
				PositionSide:           "SHORT",
				EntryPrice:             "0.0",
				MaxQty:                 "50",
			},
		},
		CanDeposit:  true,
		CanTrade:    true,
		CanWithdraw: true,
		FeeTier:     2,
		UpdateTime:  0,
	}
	s.assertAccountEqual(e, res)
}

func (s *accountServiceTestSuite) assertAccountEqual(e, a *Account) {
	r := s.r()
	r.Equal(e.CanDeposit, a.CanDeposit, "CanDeposit")
	r.Equal(e.CanTrade, a.CanTrade, "CanTrade")
	r.Equal(e.CanWithdraw, a.CanWithdraw, "CanWithdraw")
	r.Equal(e.FeeTier, a.FeeTier, "FeeTier")
	r.Equal(e.UpdateTime, a.UpdateTime, "UpdateTime")

	r.Len(a.Assets, len(e.Assets))
	for i := 0; i < len(a.Assets); i++ {
		r.Equal(e.Assets[i].Asset, a.Assets[i].Asset, "Asset")
		r.Equal(e.Assets[i].AvailableBalance, a.Assets[i].AvailableBalance, "AvailableBalance")
		r.Equal(e.Assets[i].CrossUnPnl, a.Assets[i].CrossUnPnl, "CrossUnPnl")
		r.Equal(e.Assets[i].CrossWalletBalance, a.Assets[i].CrossWalletBalance, "CrossWalletBalance")
		r.Equal(e.Assets[i].InitialMargin, a.Assets[i].InitialMargin, "InitialMargin")
		r.Equal(e.Assets[i].MaintMargin, a.Assets[i].MaintMargin, "MaintMargin")
		r.Equal(e.Assets[i].MarginBalance, a.Assets[i].MarginBalance, "MarginBalance")
		r.Equal(e.Assets[i].MaxWithdrawAmount, a.Assets[i].MaxWithdrawAmount, "MaxWithdrawAmount")
		r.Equal(e.Assets[i].OpenOrderInitialMargin, a.Assets[i].OpenOrderInitialMargin, "OpenOrderInitialMargin")
		r.Equal(e.Assets[i].PositionInitialMargin, e.Assets[i].PositionInitialMargin, "PossitionInitialMargin")
		r.Equal(e.Assets[i].UnrealizedProfit, a.Assets[i].UnrealizedProfit, "UnrealizedProfit")
		r.Equal(e.Assets[i].WalletBalance, a.Assets[i].WalletBalance, "WalletBalance")
	}

	r.Len(a.Positions, len(e.Positions))
	for i := 0; i < len(a.Positions); i++ {
		r.Equal(e.Positions[i].EntryPrice, a.Positions[i].EntryPrice, "EntryPrice")
		r.Equal(e.Positions[i].InitialMargin, a.Positions[i].InitialMargin, "InitialMargin")
		r.Equal(e.Positions[i].Isolated, a.Positions[i].Isolated, "Isolated")
		r.Equal(e.Positions[i].Leverage, a.Positions[i].Leverage, "Leverage")
		r.Equal(e.Positions[i].MaintMargin, a.Positions[i].MaintMargin, "MaintMargin")
		r.Equal(e.Positions[i].MaxQty, a.Positions[i].MaxQty, "MaxQty")
		r.Equal(e.Positions[i].OpenOrderInitialMargin, a.Positions[i].OpenOrderInitialMargin, "OpenOrderInitialMargin")
		r.Equal(e.Positions[i].PositionInitialMargin, a.Positions[i].PositionInitialMargin, "PositionInitialMargin")
		r.Equal(e.Positions[i].PositionSide, a.Positions[i].PositionSide, "PositionSide")
		r.Equal(e.Positions[i].Symbol, a.Positions[i].Symbol, "Symbol")
		r.Equal(e.Positions[i].UnrealizedProfit, a.Positions[i].UnrealizedProfit, "UnrealizedProfit")
		r.Equal(e.Positions[i].PositionAmt, a.Positions[i].PositionAmt, "PositionAmt")
	}
}
