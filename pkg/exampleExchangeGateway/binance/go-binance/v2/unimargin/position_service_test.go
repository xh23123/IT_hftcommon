package unimargin

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type positionServiceTestSuite struct {
	baseTestSuite
}

func TestPositionService(t *testing.T) {
	suite.Run(t, new(positionServiceTestSuite))
}

func (s *positionServiceTestSuite) TestChangeLeverage() {
	data := []byte(`{
		"leverage": 21,
		"maxQty": "1000",
		"symbol": "BTCUSD_200925"
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	symbol := "BTCUSD_200925"
	leverage := 21
	s.assertReq(func(r *request) {
		e := newSignedRequest().setFormParams(params{
			"symbol":   symbol,
			"leverage": leverage,
		})
		s.assertRequestEqual(e, r)
	})
	res, err := s.client.NewChangeLeverageService().Symbol(symbol).Leverage(leverage).Do(newContext())
	s.r().NoError(err)
	e := &SymbolLeverage{
		Symbol:      symbol,
		Leverage:    leverage,
		MaxQuantity: "1000",
	}
	s.r().Equal(e.Symbol, res.Symbol, "Symbol")
	s.r().Equal(e.Leverage, res.Leverage, "Leverage")
	s.r().Equal(e.MaxQuantity, res.MaxQuantity, "MaxQuantity")
}

func (s *positionServiceTestSuite) TestChangeMarginType() {
	data := []byte(`{
		"code": 200,
		"msg": "success"
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	symbol := "BTCUSDT"
	marginType := MarginTypeIsolated
	s.assertReq(func(r *request) {
		e := newSignedRequest().setFormParams(params{
			"symbol":     symbol,
			"marginType": marginType,
		})
		s.assertRequestEqual(e, r)
	})
	err := s.client.NewChangeMarginTypeService().Symbol(symbol).MarginType(marginType).Do(newContext())
	s.r().NoError(err)
}

func (s *positionServiceTestSuite) TestUpdatePositionMargin() {
	data := []byte(`{
		"amount": 100.0,
		"code": 200,
		"msg": "Successfully modify position margin.",
		"type": 1
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	symbol := "BTCUSDT"
	positionSide := PositionSideTypeLong

	amount := "100.0"
	actionType := 1
	s.assertReq(func(r *request) {
		e := newSignedRequest().setFormParams(params{
			"symbol":       symbol,
			"positionSide": positionSide,
			"amount":       amount,
			"type":         actionType,
		})
		s.assertRequestEqual(e, r)
	})
	err := s.client.NewUpdatePositionMarginService().Symbol(symbol).
		PositionSide(positionSide).Amount(amount).Type(actionType).Do(newContext())
	s.r().NoError(err)
}

func (s *positionServiceTestSuite) TestGetPositionMode() {
	data := []byte(`{
		"dualSidePosition": true
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setFormParams(params{})
		s.assertRequestEqual(e, r)
	})
	res, err := s.client.NewGetPositionModeService().Do(newContext())
	s.r().NoError(err)
	s.r().Equal(res.DualSidePosition, true)
}
