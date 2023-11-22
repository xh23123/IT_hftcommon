package unimargin

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type baseOrderTestSuite struct {
	baseTestSuite
}

type orderServiceTestSuite struct {
	baseOrderTestSuite
}

func TestOrderService(t *testing.T) {
	suite.Run(t, new(orderServiceTestSuite))
}

func (s *baseOrderTestSuite) assertCreateOrderResponseEqual(e, a *CreateOrderResponse) {
	r := s.r()
	r.Equal(e.ClientOrderID, a.ClientOrderID, "ClientOrderID")
	r.Equal(e.CumQuantity, a.CumQuantity, "CumQuantity")
	r.Equal(e.CumBase, a.CumBase, "CumBase")
	r.Equal(e.AvgPrice, a.AvgPrice, "AvgPrice")
	r.Equal(e.Pair, a.Pair, "Pair")
	r.Equal(e.OrigType, a.OrigType, "OrigType")
	r.Equal(e.PriceProtect, a.PriceProtect, "PriceProtect")
	r.Equal(e.ExecutedQuantity, a.ExecutedQuantity, "ExecutedQuantity")
	r.Equal(e.OrderID, a.OrderID, "OrderID")
	r.Equal(e.OrigQuantity, a.OrigQuantity, "OrigQuantity")
	r.Equal(e.PositionSide, a.PositionSide, "PositionSide")
	r.Equal(e.Price, a.Price, "Price")
	r.Equal(e.ReduceOnly, a.ReduceOnly, "ReduceOnly")
	r.Equal(e.Side, a.Side, "Side")
	r.Equal(e.Status, a.Status, "Status")
	r.Equal(e.StopPrice, a.StopPrice, "StopPrice")
	r.Equal(e.Symbol, a.Symbol, "Symbol")
	r.Equal(e.TimeInForce, a.TimeInForce, "TimeInForce")
	r.Equal(e.Type, a.Type, "Type")
	r.Equal(e.UpdateTime, a.UpdateTime, "UpdateTime")
	r.Equal(e.WorkingType, a.WorkingType, "WorkingType")
	r.Equal(e.ActivatePrice, a.ActivatePrice, "ActivatePrice")
	r.Equal(e.PriceRate, a.PriceRate, "PriceRate")
	r.Equal(e.ClosePosition, a.ClosePosition, "ClosePosition")
}

func (s *baseOrderTestSuite) assertOrderEqual(e, a *Order) {
	r := s.r()
	r.Equal(e.AvgPrice, a.AvgPrice, "AvgPrice")
	r.Equal(e.ClientOrderID, a.ClientOrderID, "ClientOrderID")
	r.Equal(e.CumBase, a.CumBase, "CumBase")
	r.Equal(e.ExecutedQuantity, a.ExecutedQuantity, "ExecutedQuantity")
	r.Equal(e.OrderID, a.OrderID, "OrderID")
	r.Equal(e.OrigQuantity, a.OrigQuantity, "OrigQuantity")
	r.Equal(e.OrigType, a.OrigType, "OrigType")
	r.Equal(e.Price, a.Price, "Price")
	r.Equal(e.ReduceOnly, a.ReduceOnly, "ReduceOnly")
	r.Equal(e.Side, a.Side, "Side")
	r.Equal(e.PositionSide, a.PositionSide, "PositionSide")
	r.Equal(e.Status, a.Status, "Status")
	r.Equal(e.StopPrice, a.StopPrice, "StopPrice")
	r.Equal(e.ClosePosition, a.ClosePosition, "ClosePosition")
	r.Equal(e.Symbol, a.Symbol, "Symbol")
	r.Equal(e.Pair, a.Pair, "Pair")
	r.Equal(e.Time, a.Time, "Time")
	r.Equal(e.TimeInForce, a.TimeInForce, "TimeInForce")
	r.Equal(e.Type, a.Type, "Type")
	r.Equal(e.ActivatePrice, a.ActivatePrice, "ActivatePrice")
	r.Equal(e.PriceRate, a.PriceRate, "PriceRate")
	r.Equal(e.UpdateTime, a.UpdateTime, "UpdateTime")
	r.Equal(e.WorkingType, a.WorkingType, "WorkingType")
	r.Equal(e.PriceProtect, a.PriceProtect, "PriceProtect")
}

func (s *orderServiceTestSuite) TestGetOrder() {
	data := []byte(`{
		"avgPrice": "0.0",
		"clientOrderId": "abc",
		"cumBase": "0",
		"executedQty": "0",
		"orderId": 1917641,
		"origQty": "0.40",
		"origType": "TRAILING_STOP_MARKET",
		"price": "0",
		"reduceOnly": false,
		"side": "BUY",
		"status": "NEW",
		"stopPrice": "9300",
		"closePosition": false,
		"symbol": "BTCUSD_200925",
		"pair": "BTCUSD",
		"time": 1579276756075,
		"timeInForce": "GTC",
		"type": "TRAILING_STOP_MARKET",
		"activatePrice": "9020",
		"priceRate": "0.3",
		"updateTime": 1579276756075,
		"workingType": "CONTRACT_PRICE",
		"priceProtect": false
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	symbol := "BTCUSD_200925"
	orderID := int64(1917641)
	origClientOrderID := "abc"
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParams(params{
			"symbol":            symbol,
			"orderId":           orderID,
			"origClientOrderId": origClientOrderID,
		})
		s.assertRequestEqual(e, r)
	})
	order, err := s.client.NewGetOrderService().Symbol(symbol).
		OrderID(orderID).OrigClientOrderID(origClientOrderID).Do(newContext())
	r := s.r()
	r.NoError(err)
	e := &Order{
		AvgPrice:         "0.0",
		ClientOrderID:    "abc",
		CumBase:          "0",
		ExecutedQuantity: "0",
		OrderID:          1917641,
		OrigQuantity:     "0.40",
		OrigType:         OrderTypeTrailingStopMarket,
		Price:            "0",
		ReduceOnly:       false,
		Side:             SideTypeBuy,
		Status:           OrderStatusTypeNew,
		StopPrice:        "9300",
		ClosePosition:    false,
		Symbol:           "BTCUSD_200925",
		Pair:             "BTCUSD",
		Time:             1579276756075,
		TimeInForce:      TimeInForceTypeGTC,
		Type:             OrderTypeTrailingStopMarket,
		ActivatePrice:    "9020",
		PriceRate:        "0.3",
		UpdateTime:       1579276756075,
		WorkingType:      WorkingTypeContractPrice,
		PriceProtect:     false,
	}
	s.assertOrderEqual(e, order)
}

func (s *orderServiceTestSuite) TestListOrders() {
	data := []byte(`[
		{
			"avgPrice": "0.0",
			"clientOrderId": "abc",
			"cumBase": "0",
			"executedQty": "0",
			"orderId": 1917641,
			"origQty": "0.40",
			"origType": "TRAILING_STOP_MARKET",
			"price": "0",
			"reduceOnly": false,
			"side": "BUY",
			"positionSide": "SHORT",
			"status": "NEW",
			"stopPrice": "9300",
			"closePosition": false,
			"symbol": "BTCUSD_200925",
			"pair": "BTCUSD",
			"time": 1579276756075,
			"timeInForce": "GTC",
			"type": "TRAILING_STOP_MARKET",
			"activatePrice": "9020",
			"priceRate": "0.3",
			"updateTime": 1579276756075,
			"workingType": "CONTRACT_PRICE",
			"priceProtect": false
		}
	  ]`)
	s.mockDo(data, nil)
	defer s.assertDo()
	symbol := "BTCUSD_200925"
	orderID := int64(1917641)
	limit := 3
	startTime := int64(1499827319559)
	endTime := int64(1499827319560)
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParams(params{
			"symbol":    symbol,
			"orderId":   orderID,
			"startTime": startTime,
			"endTime":   endTime,
			"limit":     limit,
		})
		s.assertRequestEqual(e, r)
	})

	orders, err := s.client.NewListOrdersService().Symbol(symbol).
		OrderID(orderID).StartTime(startTime).EndTime(endTime).
		Limit(limit).Do(newContext())
	r := s.r()
	r.NoError(err)
	r.Len(orders, 1)
	e := &Order{
		AvgPrice:         "0.0",
		ClientOrderID:    "abc",
		CumBase:          "0",
		ExecutedQuantity: "0",
		OrderID:          1917641,
		OrigQuantity:     "0.40",
		OrigType:         OrderTypeTrailingStopMarket,
		Price:            "0",
		ReduceOnly:       false,
		Side:             SideTypeBuy,
		PositionSide:     PositionSideTypeShort,
		Status:           OrderStatusTypeNew,
		StopPrice:        "9300",
		ClosePosition:    false,
		Symbol:           "BTCUSD_200925",
		Pair:             "BTCUSD",
		Time:             1579276756075,
		TimeInForce:      TimeInForceTypeGTC,
		Type:             OrderTypeTrailingStopMarket,
		ActivatePrice:    "9020",
		PriceRate:        "0.3",
		UpdateTime:       1579276756075,
		WorkingType:      WorkingTypeContractPrice,
		PriceProtect:     false,
	}
	s.assertOrderEqual(e, orders[0])
}

func (s *orderServiceTestSuite) assertCancelOrderResponseEqual(e, a *CancelOrderResponse) {
	r := s.r()
	r.Equal(e.AvgPrice, a.AvgPrice, "AvgPrice")
	r.Equal(e.ClientOrderID, a.ClientOrderID, "ClientOrderID")
	r.Equal(e.CumQuantity, a.CumQuantity, "CumQuantity")
	r.Equal(e.CumBase, a.CumBase, "CumBase")
	r.Equal(e.ExecutedQuantity, a.ExecutedQuantity, "ExecutedQuantity")
	r.Equal(e.OrderID, a.OrderID, "OrderID")
	r.Equal(e.OrigQuantity, a.OrigQuantity, "OrigQuantity")
	r.Equal(e.OrigType, a.OrigType, "OrigType")
	r.Equal(e.Price, a.Price, "Price")
	r.Equal(e.ReduceOnly, a.ReduceOnly, "ReduceOnly")
	r.Equal(e.Side, a.Side, "Side")
	r.Equal(e.PositionSide, a.PositionSide, "PositionSide")
	r.Equal(e.Status, a.Status, "Status")
	r.Equal(e.StopPrice, a.StopPrice, "StopPrice")
	r.Equal(e.ClosePosition, a.ClosePosition, "ClosePosition")
	r.Equal(e.Symbol, a.Symbol, "Symbol")
	r.Equal(e.Pair, a.Pair, "Pair")
	r.Equal(e.TimeInForce, a.TimeInForce, "TimeInForce")
	r.Equal(e.Type, a.Type, "Type")
	r.Equal(e.ActivatePrice, a.ActivatePrice, "ActivatePrice")
	r.Equal(e.PriceRate, a.PriceRate, "PriceRate")
	r.Equal(e.UpdateTime, a.UpdateTime, "UpdateTime")
	r.Equal(e.WorkingType, a.WorkingType, "WorkingType")
	r.Equal(e.PriceProtect, a.PriceProtect, "PriceProtect")
}

func (s *orderServiceTestSuite) TestListLiquidationOrders() {
	data := []byte(`[
		{
			"symbol": "BTCUSD_200925",
			"price": "9425.5",
			"origQty": "1",
			"executedQty": "1",
			"avragePrice": "9496.5",
			"status": "FILLED",
			"timeInForce": "IOC",
			"type": "LIMIT",
			"side": "SELL",
			"time": 1591154240949
		}
	]`)
	s.mockDo(data, nil)
	defer s.assertDo()

	symbol := "BTCUSD_200925"
	startTime := int64(1568014460893)
	endTime := int64(1568014460894)
	limit := 1
	s.assertReq(func(r *request) {
		e := newRequest().setParams(params{
			"symbol":    symbol,
			"startTime": startTime,
			"endTime":   endTime,
			"limit":     limit,
		})
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewListLiquidationOrdersService().Symbol(symbol).
		StartTime(startTime).EndTime(endTime).Limit(limit).Do(newContext())
	r := s.r()
	r.NoError(err)
	e := []*LiquidationOrder{
		{
			Symbol:           symbol,
			Price:            "9425.5",
			OrigQuantity:     "1",
			ExecutedQuantity: "1",
			AveragePrice:     "9496.5",
			Status:           OrderStatusTypeFilled,
			TimeInForce:      TimeInForceTypeIOC,
			Type:             OrderTypeLimit,
			Side:             SideTypeSell,
			Time:             1591154240949,
		},
	}
	s.r().Len(res, len(e))
	for i := range res {
		s.assertLiquidationEqual(e[i], res[i])
	}
}

func (s *orderServiceTestSuite) assertLiquidationEqual(e, a *LiquidationOrder) {
	r := s.r()
	r.Equal(e.Symbol, a.Symbol, "Symbol")
	r.Equal(e.Price, a.Price, "Price")
	r.Equal(e.OrigQuantity, a.OrigQuantity, "OrigQuantity")
	r.Equal(e.ExecutedQuantity, a.ExecutedQuantity, "ExecutedQuantity")
	r.Equal(e.AveragePrice, a.AveragePrice, "AveragePrice")
	r.Equal(e.Status, a.Status, "Status")
	r.Equal(e.TimeInForce, a.TimeInForce, "TimeInForce")
	r.Equal(e.Type, a.Type, "Type")
	r.Equal(e.Side, a.Side, "Side")
	r.Equal(e.Time, a.Time, "Time")
}
