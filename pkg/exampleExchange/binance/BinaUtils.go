package bina

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/xh23123/IT_hftcommon/pkg/common"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchange/binance/go-binance/v2"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchange/binance/go-binance/v2/futures"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchange/binance/go-binance/v2/unimargin"
	"go.uber.org/zap"
)

var binanceSysToExMap cmap.ConcurrentMap
var binanceExToSysMap cmap.ConcurrentMap

const BinanceExchangeStr = "binance:"

type AccountInfo struct {
	Name  string
	Index common.AccountIdx
}

func init() {
	binanceSysToExMap = cmap.New()
	binanceExToSysMap = cmap.New()
}

func RegisterSystemSymbols(symbols []string) {
	for _, symbol := range symbols {
		GetExchangeSymbol(symbol)
	}
}

func convertSysToExchangeSymbol(sysSymbol string) string {

	if !strings.Contains(sysSymbol, "-") {
		panic("invalid sysSymbol " + sysSymbol)
	}

	return strings.Replace(sysSymbol, "-", "", -1)
}

func GetExchangeSymbol(sysSymbol string) string {
	if v, ok := binanceSysToExMap.Get(sysSymbol); ok {
		return v.(string)
	} else {
		exchangeSymbol := convertSysToExchangeSymbol(sysSymbol)
		binanceSysToExMap.Set(sysSymbol, exchangeSymbol)
		binanceExToSysMap.Set(exchangeSymbol, sysSymbol)
		return exchangeSymbol
	}
}

func GetCoinFutureExchangeSymbol(sysSymbol string) string {
	return fmt.Sprintf("%s_PERP", GetExchangeSymbol(sysSymbol))
}

func GetCoinFutureSysSymbol(exchangeSymbol string) string {
	if symbol, err := TryGetCoinFutureSysSymbol(exchangeSymbol); err == nil {
		return symbol
	} else {
		panic(" GetCoinFutureSysSymbol cant find " + exchangeSymbol)
	}
}

func TryGetCoinFutureSysSymbol(exchangeSymbol string) (string, error) {
	exchangeSymbol = exchangeSymbol[:len(exchangeSymbol)-5]
	if v, ok := binanceExToSysMap.Get(exchangeSymbol); ok {
		return v.(string), nil
	} else {
		return "", errors.New("TryGetCoinFutureSysSymbol Not sysmbol found " + exchangeSymbol)
	}
}

func GetSysSymbol(exchangeSymbol string) string {
	if symbol, err := TryGetSysSymbol(exchangeSymbol); err == nil {
		return symbol
	} else {
		panic(" GetSysSymbol cant find " + exchangeSymbol)
	}
}

func TryGetSysSymbol(exchangeSymbol string) (string, error) {
	if v, ok := binanceExToSysMap.Get(exchangeSymbol); ok {
		return v.(string), nil
	} else {
		return "", errors.New("Not sysmbol found " + exchangeSymbol)
	}
}

func convertCoinFutureToExchangeSymbolsSlice(sysSymbols []string) []string {
	exchangeSymbols := make([]string, 0, len(sysSymbols))

	for k := range sysSymbols {
		exchangeSymbols = append(exchangeSymbols, GetCoinFutureExchangeSymbol(sysSymbols[k]))
	}

	return exchangeSymbols
}

func ConvertToExchangeSymbolsSlice(sysSymbols []string) []string {
	exchangeSymbols := make([]string, 0, len(sysSymbols))

	for k := range sysSymbols {
		exchangeSymbols = append(exchangeSymbols, GetExchangeSymbol(sysSymbols[k]))
	}

	return exchangeSymbols
}

func ConvertToExchangeSymbolsMap(sysSymbols map[string]string) map[string]string {
	exchangeSymbols := map[string]string{}

	for k := range sysSymbols {
		exchangeSymbols[GetExchangeSymbol(k)] = sysSymbols[k]
	}

	return exchangeSymbols
}

func GetBinanceAccount(accountIndex common.AccountIdx) string {
	return fmt.Sprintf("binance:%d", accountIndex)
}

type ConcurrentTickTriggerInfoMap struct {
	TickTrigRecord map[common.TransactionID]map[string]*common.TickTrigInfo
	mux            sync.Mutex
}

func (c *ConcurrentTickTriggerInfoMap) GetTriggerInfo(transactionId common.TransactionID, symbol string) (*common.TickTrigInfo, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.TickTrigRecord[transactionId]; !ok {
		c.TickTrigRecord[transactionId] = make(map[string]*common.TickTrigInfo)
	}

	if triggerInfo, ok := c.TickTrigRecord[transactionId][symbol]; ok {
		return triggerInfo, true
	} else {
		trickInfo := &common.TickTrigInfo{}
		c.TickTrigRecord[transactionId][symbol] = trickInfo
		return trickInfo, false
	}
}

type ConcurrentCancelMap struct {
	CancelCtx map[common.ResetID]chan struct{}
	mux       sync.Mutex
}

func (c *ConcurrentCancelMap) SetCancelChannel(resetId common.ResetID, channel chan struct{}) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if stopC, ok := c.CancelCtx[resetId]; ok {
		common.Logger.Warn("SetCancelChannel already has stopC channel: " + string(resetId))
		close(stopC)
	}
	c.CancelCtx[resetId] = channel
}

func (c *ConcurrentCancelMap) CloseCancelChannels(rs []common.ResetID) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, r := range rs {
		if stopC, ok := c.CancelCtx[r]; ok {
			if stopC != nil {
				close(stopC)
			}

			delete(c.CancelCtx, r)
		}
	}
}

func TypeConvert(src interface{}, targ interface{}) error {
	resStr, err := json.Marshal(src)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(resStr, targ); err != nil {
		return err
	}
	return nil
}

func GetFutureBalance(client *unimargin.Client) (common.WsFutureBalance, error) {
	res, err := client.NewGetBalanceService().Do(context.Background())
	if err != nil {
		common.Logger.Error("GetFutureBalance failed ",
			zap.String("err", err.Error()))

		return nil, err
	}
	balances := common.WsFutureBalance{}
	for _, r := range res {
		balances[r.Asset] = &common.Balance{
			Balance:          common.Str2Float(r.UmWalletBalance),
			MarginBalance:    common.Str2Float(r.UmWalletBalance) + common.Str2Float(r.UmUnrealizedPNL),
			AvailableBalance: common.Str2Float(r.TotalWalletBalance),
		}
	}

	return balances, nil
}

func GetCoinFutureBalance(client *unimargin.Client) (common.WsFutureBalance, error) {
	res, err := client.NewGetBalanceService().Do(context.Background())
	if err != nil {
		common.Logger.Error("GetCoinFutureBalance failed ",
			zap.String("err", err.Error()))

		return nil, err
	}
	balances := common.WsFutureBalance{}
	for _, r := range res {
		balances[r.Asset] = &common.Balance{
			Balance:          common.Str2Float(r.CmWalletBalance),
			MarginBalance:    common.Str2Float(r.CmWalletBalance) + common.Str2Float(r.CmUnrealizedPNL),
			AvailableBalance: common.Str2Float(r.TotalWalletBalance),
		}
	}

	return balances, nil
}

func GetSpotBalance(client *binance.Client) (common.SpotBalance, error) {
	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	balances := common.SpotBalance{}
	for _, r := range res.Balances {
		free, _ := strconv.ParseFloat(r.Free, 64)
		lock, _ := strconv.ParseFloat(r.Locked, 64)
		balance := common.Balance{Balance: free + lock}
		balances[r.Asset] = &balance
	}
	return balances, nil

}

func GetPremiumIndex(futureCli *futures.Client, symbol string) []*common.PremiumIndexInfo {
	if preSlice, err := futureCli.NewPremiumIndexService().Symbol(symbol).Do(context.Background()); err == nil {
		ret := make([]*common.PremiumIndexInfo, 0, len(preSlice))

		for _, v := range preSlice {
			cPremium := common.PremiumIndexInfo{}
			cPremium.Symbol = GetSysSymbol(v.Symbol)
			cPremium.MarkPrice = v.MarkPrice
			cPremium.LastFundingRate = v.LastFundingRate
			cPremium.NextFundingTime = v.NextFundingTime
			cPremium.Time = v.Time
			ret = append(ret, &cPremium)
		}
		return ret
	} else {
		common.Logger.Error("Binance GetPremiumIndex Error", zap.String("error", err.Error()))
		return nil
	}
}

func SetMultiAssetMargin(futureCli *futures.Client, MultiAssetMargin bool) {
	err := futureCli.NewChangeMultiAssetsMarginService().MultiAssetMargin(MultiAssetMargin).Do(context.Background())
	if err != nil {
		common.Logger.Error("Binance SetMultiAssetMargin Error",
			zap.String("error", err.Error()),
			zap.String("APIKey", futureCli.APIKey),
			zap.Bool("MultiAssetMargin", MultiAssetMargin))
	}
}
