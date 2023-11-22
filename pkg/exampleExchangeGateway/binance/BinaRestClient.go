package bina

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/xh23123/IT_hftcommon/pkg/common"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2/delivery"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2/futures"
	"go.uber.org/zap"
)

var _ common.RestClientInterface = (*BinaRestClient)(nil)

type BinaRestClient struct {
	spotCli       *binance.Client
	futureCli     *futures.Client
	coinFutureCli *delivery.Client
}

func NewRestClient(config map[string]string) *BinaRestClient {
	apiKey := config["apiKey"]
	secretKey := config["secretKey"]
	client := BinaRestClient{}
	client.spotCli = binance.NewClient(apiKey, secretKey)
	client.futureCli = futures.NewClient(apiKey, secretKey)
	client.coinFutureCli = delivery.NewClient(apiKey, secretKey)
	return &client
}

func (r *BinaRestClient) GetPremiumIndex(symbol string) []*common.PremiumIndexInfo {
	return GetPremiumIndex(r.futureCli, GetExchangeSymbol(symbol))
}

func (r *BinaRestClient) SetMultiAssetMargin(MultiAssetMargin bool) {
	SetMultiAssetMargin(r.futureCli, MultiAssetMargin)
}

func (r *BinaRestClient) GetOrder(symbol string, transactionId common.TransactionID, origClientOrderID string) *common.Order {
	openOrders := r.GetOrders(symbol, transactionId)

	for _, o := range openOrders {
		if origClientOrderID == o.Cid {
			return o
		}
	}
	return nil
}

func (r *BinaRestClient) GetOrders(symbol string, transactionId common.TransactionID) []*common.Order {
	switch transactionId {
	case common.SpotID:
		{
			return r.getSpotOrders(symbol)
		}
	case common.FutureID:
		{
			return r.getFutureOrders(symbol)
		}
	case common.MarginID:
		{
			return r.getMarginOrders(symbol)
		}
	case common.CoinFutureID:
		{
			return r.getCoinFutureOrders(symbol)
		}
	default:
		panic("BinaRestClient::GetOrders invalid transactionId : " + transactionId)
	}
}

func (r *BinaRestClient) getSpotOrders(symbol string) []*common.Order {
	restService := r.spotCli.NewListOpenOrdersService()
	restService.Symbol(GetExchangeSymbol(symbol))
	openOrders, err := restService.Do(context.Background())
	if err != nil {
		common.Logger.Error("Binance GetOrders Error", zap.String("error", err.Error()))
		return nil
	}
	ret := make([]*common.Order, 0, len(openOrders))
	for _, o := range openOrders {
		ret = append(ret,
			&common.Order{
				Exchange:    common.BINANCEID,
				Transaction: common.SpotID,
				Symbol:      GetSysSymbol(o.Symbol),
				Id:          common.Int2Str(o.OrderID),
				Cid:         o.ClientOrderID,
				Side:        common.SideID(o.Side),
				Type:        common.OrderTypeID(o.Type),
				FilledSize:  common.Str2Float(o.ExecutedQuantity),
				Size:        common.Str2Float(o.OrigQuantity),
				Price:       common.Str2Float(o.Price),
				CreateTime:  o.Time,
				Status:      common.Normal})
	}

	return ret
}

func (r *BinaRestClient) getMarginOrders(symbol string) []*common.Order {
	restService := r.spotCli.NewListMarginOpenOrdersService()
	restService.Symbol(GetExchangeSymbol(symbol))
	openOrders, err := restService.Do(context.Background())
	if err != nil {
		common.Logger.Error("Binance getMarginOrders Error", zap.String("error", err.Error()))
		return nil
	}
	ret := make([]*common.Order, 0, len(openOrders))
	for _, o := range openOrders {
		ret = append(ret,
			&common.Order{
				Exchange:    common.BINANCEID,
				Transaction: common.MarginID,
				Symbol:      GetSysSymbol(o.Symbol),
				Id:          common.Int2Str(o.OrderID),
				Cid:         o.ClientOrderID,
				Side:        common.SideID(o.Side),
				Type:        common.OrderTypeID(o.Type),
				IsIsolated:  o.IsIsolated,
				FilledSize:  common.Str2Float(o.ExecutedQuantity),
				Size:        common.Str2Float(o.OrigQuantity),
				Price:       common.Str2Float(o.Price),
				CreateTime:  o.Time,
				Status:      common.Normal})
	}

	return ret
}

func (r *BinaRestClient) getCoinFutureOrders(symbol string) []*common.Order {
	restService := r.coinFutureCli.NewListOpenOrdersService()
	restService.Symbol(GetCoinFutureExchangeSymbol(symbol))
	openOrders, err := restService.Do(context.Background())
	if err != nil {
		common.Logger.Error("Binance getCoinFutureOrders Error", zap.String("error", err.Error()))
		return nil
	}
	ret := make([]*common.Order, 0, len(openOrders))
	for _, o := range openOrders {
		ret = append(ret,
			&common.Order{
				Exchange:    common.BINANCEID,
				Transaction: common.CoinFutureID,
				Symbol:      GetCoinFutureSysSymbol(o.Symbol),
				Id:          common.Int2Str(o.OrderID),
				Cid:         o.ClientOrderID,
				Side:        common.SideID(o.Side),
				Type:        common.OrderTypeID(o.Type),
				FilledSize:  common.Str2Float(o.ExecutedQuantity),
				Size:        common.Str2Float(o.OrigQuantity),
				Price:       common.Str2Float(o.Price),
				CreateTime:  o.Time,
				Status:      common.Normal})
	}

	return ret
}

func (r *BinaRestClient) getFutureOrders(symbol string) []*common.Order {
	restService := r.futureCli.NewListOpenOrdersService()
	restService.Symbol(GetExchangeSymbol(symbol))
	openOrders, err := restService.Do(context.Background())
	if err != nil {
		common.Logger.Error("Binance GetFutureOrder Error", zap.String("error", err.Error()))
		return nil
	}

	ret := make([]*common.Order, 0, len(openOrders))

	for _, o := range openOrders {
		ret = append(ret,
			&common.Order{
				Exchange:    common.BINANCEID,
				Transaction: common.FutureID,
				Symbol:      GetSysSymbol(o.Symbol),
				Id:          common.Int2Str(o.OrderID),
				Cid:         o.ClientOrderID,
				Side:        common.SideID(o.Side),
				Type:        common.OrderTypeID(o.Type),
				FilledSize:  common.Str2Float(o.ExecutedQuantity),
				Size:        common.Str2Float(o.OrigQuantity),
				Price:       common.Str2Float(o.Price),
				CreateTime:  o.Time,
				Status:      common.Normal})
	}
	return ret
}

func (r *BinaRestClient) GetSpotBalance() (common.SpotBalance, error) {
	return GetSpotBalance(r.spotCli)
}

func (r *BinaRestClient) GetMarginBalance() (common.MarginBalance, error) {
	ret := common.MarginBalance{}
	response, err := r.spotCli.NewGetMarginAccountService().Do(context.Background())
	if err != nil {
		return ret, err
	}
	ret.MarginLevel = common.Str2Float(response.MarginLevel)
	for _, r := range response.UserAssets {
		userAsset := common.UserAsset{}
		userAsset.Asset = r.Asset
		userAsset.Borrowed = common.Str2Float(r.Borrowed)
		userAsset.Free = common.Str2Float(r.Free)
		userAsset.Interest = common.Str2Float(r.Interest)
		userAsset.Locked = common.Str2Float(r.Locked)
		userAsset.NetAsset = common.Str2Float(r.NetAsset)
		ret.UserAssets = append(ret.UserAssets, userAsset)
	}
	return ret, nil
}

func (r *BinaRestClient) GetFutureBalancePosition() (common.WsFutureBalance, common.WsFuturePosition, error) {
	res, err := r.futureCli.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, nil, err
	}
	balances := common.WsFutureBalance{}
	for _, r := range res.Assets {

		balances[r.Asset] = &common.Balance{
			Balance:          common.Str2Float(r.WalletBalance),
			MarginBalance:    common.Str2Float(r.MarginBalance),
			AvailableBalance: common.Str2Float(r.AvailableBalance),
		}
	}
	positions := common.WsFuturePosition{}
	for _, r := range res.Positions {
		symbol, err := TryGetSysSymbol(r.Symbol)
		if err != nil {
			continue
		}
		pos, ok := positions[symbol]
		if !ok {
			positions[symbol] = &common.FuturePosition{}
			pos = positions[symbol]
		}
		if r.PositionSide == futures.PositionSideTypeLong {
			pos.LONG = &common.SidePosition{
				Amount:        common.Str2Float(r.PositionAmt),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedProfit),
			}
		} else if r.PositionSide == futures.PositionSideTypeShort {
			pos.SHORT = &common.SidePosition{
				Amount:        common.Str2Float(r.PositionAmt),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedProfit),
			}
		} else if r.PositionSide == futures.PositionSideTypeBoth {
			pos.BOTH = &common.SidePosition{
				Amount:        common.Str2Float(r.PositionAmt),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedProfit),
			}
		}
	}
	return balances, positions, nil
}

func (r *BinaRestClient) GetFutureKlines(symbol string, interval common.IntervalID, limit int, startTime int64, endTime int64) ([]*common.Kline, error) {
	service := r.futureCli.NewKlinesService().Symbol(GetExchangeSymbol(symbol)).
		Interval(string(interval)).Limit(limit)

	if startTime != 0 {
		service = service.StartTime(startTime)
	}
	if endTime != 0 {
		service = service.EndTime(endTime)
	}
	if klines, err := service.Do(context.Background()); err != nil {
		return nil, err
	} else {
		retKlines := make([]*common.Kline, 0, len(klines))

		for _, v := range klines {
			retKlines = append(retKlines, &common.Kline{
				OpenTime:         v.OpenTime,
				Open:             v.Open,
				High:             v.High,
				Low:              v.Low,
				Close:            v.Close,
				Volume:           v.Volume,
				QuoteAssetVolume: v.QuoteAssetVolume,
			})
		}

		return retKlines, nil
	}
}

func (r *BinaRestClient) GetSpotKlines(symbol string, interval common.IntervalID, limit int, startTime int64, endTime int64) ([]*common.Kline, error) {

	service := r.spotCli.NewKlinesService().Symbol(GetExchangeSymbol(symbol)).
		Interval(string(interval)).Limit(limit)

	if startTime != 0 {
		service = service.StartTime(startTime)
	}
	if endTime != 0 {
		service = service.EndTime(endTime)
	}
	if klines, err := service.Do(context.Background()); err != nil {
		return nil, err
	} else {
		retKlines := make([]*common.Kline, 0, len(klines))

		for _, v := range klines {
			retKlines = append(retKlines, &common.Kline{
				OpenTime:         v.OpenTime,
				Open:             v.Open,
				High:             v.High,
				Low:              v.Low,
				Close:            v.Close,
				Volume:           v.Volume,
				QuoteAssetVolume: v.QuoteAssetVolume,
			})
		}

		return retKlines, nil
	}
}

func (r *BinaRestClient) MarginLoan(asset string,
	isIsolated bool,
	symbol string,
	amount float64) (*common.TransactionResponse, error) {
	binanceRes, err := r.spotCli.NewMarginLoanService().Asset(asset).IsIsolated(isIsolated).Symbol(symbol).
		Amount(common.Float2Str(amount)).Do(context.Background())

	if err != nil {
		return nil, err
	} else {
		ret := common.TransactionResponse{}
		if err = TypeConvert(*binanceRes, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

func (r *BinaRestClient) MarginRepay(asset string,
	isIsolated bool,
	symbol string,
	amount float64) (*common.TransactionResponse, error) {
	binanceRes, err := r.spotCli.NewMarginRepayService().Asset(asset).IsIsolated(isIsolated).Symbol(symbol).
		Amount(common.Float2Str(amount)).Do(context.Background())

	if err != nil {
		return nil, err
	} else {
		ret := common.TransactionResponse{}
		if err = TypeConvert(*binanceRes, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

func (r *BinaRestClient) MarginAllAssets() ([]*common.MarginAsset, error) {
	binanceRes, err := r.spotCli.NewGetAllMarginAssetsService().Do(context.Background())

	if err != nil {
		return nil, err
	} else {
		ret := []*common.MarginAsset{}
		if err = TypeConvert(binanceRes, &ret); err != nil {
			return nil, err
		}
		return ret, nil
	}
}

func (r *BinaRestClient) MarginAllPairs() ([]*common.MarginAllPair, error) {
	binanceRes, err := r.spotCli.NewGetMarginAllPairsService().Do(context.Background())

	if err != nil {
		return nil, err
	} else {
		ret := []*common.MarginAllPair{}
		if err = TypeConvert(binanceRes, &ret); err != nil {
			return nil, err
		}
		return ret, nil
	}
}

func (r *BinaRestClient) CrossMarginCollateralRatio() (res []*common.CrossMarginCollateralRatio, err error) {
	binanceRes, err := r.spotCli.NewGetCrossMarginCollateralRatioService().Do(context.Background())

	if err != nil {
		return nil, err
	} else {
		ret := []*common.CrossMarginCollateralRatio{}
		if err = TypeConvert(binanceRes, &ret); err != nil {
			return nil, err
		}
		return ret, nil
	}
}

func (r *BinaRestClient) NextHourlyInterestRates(assets []string, isIsolated bool) ([]*common.NextHourlyInterestRate, error) {
	binanceRes, err := r.spotCli.NewGetNextHourlyInterestRateService().
		Assets(assets).
		IsIsolated(isIsolated).
		Do(context.Background())

	if err != nil {
		return nil, err
	} else {
		ret := []*common.NextHourlyInterestRate{}
		if err = TypeConvert(binanceRes, &ret); err != nil {
			return nil, err
		}
		return ret, nil
	}
}

func (r *BinaRestClient) GetDustAssets() (*common.ListDustResponse, error) {
	if binanceRes, err := r.spotCli.NewListDustService().Do(context.Background()); err != nil {
		return nil, err
	} else {
		ret := common.ListDustResponse{}
		if err = TypeConvert(*binanceRes, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

func (r *BinaRestClient) ConvertDustAssets(assets []string) (*common.DustTransferResponse, error) {

	if binanceRes, err := r.spotCli.NewDustTransferService().Asset(assets).Do(context.Background()); err != nil {
		return nil, err
	} else {
		ret := common.DustTransferResponse{}
		if err = TypeConvert(*binanceRes, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

func (r *BinaRestClient) GetMarginDustAssets() (*[]common.ListMarginDustResponse, error) {
	if binanceRes, err := r.spotCli.NewListMarginDustService().Do(context.Background()); err != nil {
		return nil, err
	} else {
		ret := []common.ListMarginDustResponse{}
		if err = TypeConvert(*binanceRes, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

func (r *BinaRestClient) ConvertMarginDustAssets(assets []string) (*[]common.ListMarginDustResponse, error) {
	if binanceRes, err := r.spotCli.NewMarginDustTransferService().Assets(assets).Do(context.Background()); err != nil {
		return nil, err
	} else {
		ret := []common.ListMarginDustResponse{}
		if err = TypeConvert(*binanceRes, &ret); err != nil {
			return nil, err
		}
		return &ret, nil
	}
}

func (r *BinaRestClient) GetCoinFutureBalancePosition() (common.WsFutureBalance, common.WsFuturePosition, error) {

	res, err := r.coinFutureCli.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, nil, err
	}
	balances := common.WsFutureBalance{}
	for _, r := range res.Assets {

		balances[r.Asset] = &common.Balance{
			Balance:          common.Str2Float(r.WalletBalance),
			MarginBalance:    common.Str2Float(r.MarginBalance),
			AvailableBalance: common.Str2Float(r.AvailableBalance),
		}
	}
	positions := common.WsFuturePosition{}
	for _, r := range res.Positions {
		symbol, err := TryGetCoinFutureSysSymbol(r.Symbol)
		if err != nil {
			continue
		}
		pos, ok := positions[symbol]
		if !ok {
			positions[symbol] = &common.FuturePosition{}
			pos = positions[symbol]
		}
		if r.PositionSide == string(delivery.PositionSideTypeLong) {
			pos.LONG = &common.SidePosition{
				Amount:        common.Str2Float(r.PositionAmt),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedProfit),
			}
		} else if r.PositionSide == string(delivery.PositionSideTypeShort) {
			pos.SHORT = &common.SidePosition{
				Amount:        common.Str2Float(r.PositionAmt),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedProfit),
			}
		} else if r.PositionSide == string(delivery.PositionSideTypeBoth) {
			pos.BOTH = &common.SidePosition{
				Amount:        common.Str2Float(r.PositionAmt),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedProfit),
			}
		}
	}
	return balances, positions, nil
}

func (r *BinaRestClient) GetExchangeInfos(transactionId common.TransactionID) (*common.ExchangeInfo, error) {

	switch transactionId {
	case common.SpotID:
		{
			res, err := r.spotCli.NewExchangeInfoService().Do(context.Background())
			marshalStr, _ := json.Marshal(res)
			newRes := common.ExchangeInfo{}
			err = json.Unmarshal(marshalStr, &newRes)
			for k := range newRes.Symbols {
				if symbol, error := TryGetSysSymbol(newRes.Symbols[k].Symbol); error == nil {
					newRes.Symbols[k].Symbol = symbol
				}

			}

			return &newRes, err
		}
	case common.FutureID:
		{
			res, _ := r.futureCli.NewExchangeInfoService().Do(context.Background())
			marshalStr, _ := json.Marshal(res)
			newRes := common.ExchangeInfo{}
			err := json.Unmarshal(marshalStr, &newRes)

			for k := range newRes.Symbols {
				if symbol, error := TryGetSysSymbol(newRes.Symbols[k].Symbol); error == nil {
					newRes.Symbols[k].Symbol = symbol
				}

			}
			return &newRes, err
		}

	default:
		panic("BinaRestClient::GetExchangeInfos invalid transactionId : " + transactionId)
	}

}

func (r *BinaRestClient) GetSuggestGasPrice() (*big.Int, error) {
	return nil, fmt.Errorf("no implement")
}
