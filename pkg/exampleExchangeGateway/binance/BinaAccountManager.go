package bina

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/xh23123/IT_hftcommon/pkg/common"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2"
	"github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance/go-binance/v2/futures"
	"go.uber.org/zap"

	cmap "github.com/orcaman/concurrent-map"
)

var _ common.AccountManagerInterface = (*BinaAccountManager)(nil)

const Rest_update_interval = 50 * time.Millisecond

type BinaAccountManager struct {
	systemAgent  common.TradeSystemAgent
	accountIndex common.AccountIdx

	spotCli      *binance.Client
	futureCli    *futures.Client
	orderManager map[common.TransactionID]common.OrderManagerInterface

	balanceManager common.BalanceManagerInterface
}

func NewBinaAccountManager(systemAgent common.TradeSystemAgent) {

	accountConfig := GetBinanceAccountConfigs(systemAgent)

	for _, v := range accountConfig {

		bam := &BinaAccountManager{
			accountIndex: v.Index,
			orderManager: map[common.TransactionID]common.OrderManagerInterface{
				common.SpotID:   systemAgent.NewOrderManager(common.BINANCEID, v.Index, common.SpotID),
				common.FutureID: systemAgent.NewOrderManager(common.BINANCEID, v.Index, common.FutureID),
			},
			balanceManager: systemAgent.NewBalanceManager(common.BINANCEID, v.Index),
		}

		apiKey := systemAgent.Config().Section(v.Name).Key("api_key").String()
		secretKey := systemAgent.Config().Section(v.Name).Key("secret_key").String()
		bam.spotCli = binance.NewClient(apiKey, secretKey)
		bam.futureCli = futures.NewClient(apiKey, secretKey)

		systemAgent.RegisterAccountManager(common.BINANCEID, bam)

		bam.initAccountInfo()
		go bam.onTimer()

	}

}

func GetBinanceAccountConfigs(systemAgent common.TradeSystemAgent) []AccountInfo {
	//pubPort, err := systemAgent.Config().Section("global").Key("pubPort").Int64()
	accounts := []AccountInfo{}

	sections := systemAgent.Config().Sections()
	expectAccountIdx := common.AccountIdx(0)
	for _, v := range sections {
		if strings.HasPrefix(v.Name(), BinanceExchangeStr) {
			expectAccountId := GetBinanceAccount(expectAccountIdx)
			if expectAccountId != v.Name() {
				panic("binance account expect account: " + expectAccountId)
			}

			accounts = append(accounts, AccountInfo{Name: v.Name(), Index: expectAccountIdx})

			expectAccountIdx += 1
		}
	}

	return accounts
}

func (b *BinaAccountManager) Process(event common.ActionEvent) {
	switch event.Action {
	case common.CREATE_LIMITTYPE_SPOT_ORDER:
		b.orderManager[common.SpotID].CreateOrderProcess(event, b.createLimitTypeSpotOrderProcess)
	case common.CREATE_MARKET_SPOT_ORDER:
		b.orderManager[common.SpotID].CreateOrderProcess(event, b.createMarketSpotOrderProcess)
	case common.CANCEL_SPOT_ORDER:
		b.orderManager[common.SpotID].CancelOrderProcess(event, b.cancelSpotOrderProcess)
	case common.CANCEL_ALL_SPOT_ORDER:
		b.orderManager[common.SpotID].CancelAllOrderProcess(event, b.cancelAllSpotOrderProcess)
	case common.CREATE_LIMITTYPE_FUTURE_ORDER:
		b.orderManager[common.FutureID].CreateOrderProcess(event, b.createLimitTypeFutureOrderProcess)
	case common.CREATE_MARKET_FUTURE_ORDER:
		b.orderManager[common.FutureID].CreateOrderProcess(event, b.createMarketFutureOrderProcess)
	case common.CANCEL_FUTURE_ORDER:
		b.orderManager[common.FutureID].CancelOrderProcess(event, b.cancelFutureOrderProcess)
	case common.CANCEL_ALL_FUTURE_ORDER:
		b.orderManager[common.FutureID].CancelAllOrderProcess(event, b.cancelFutureAllOrderProcess)
	default:
		common.Logger.Error("BinaAccountManager::Process unknown action ",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.Any("actionID", event.Action))
	}
}

func (b *BinaAccountManager) GetAllOrders(transactionId common.TransactionID) cmap.ConcurrentMap {
	return b.orderManager[transactionId].AllOpenOrders()
}
func (b *BinaAccountManager) GetBalance(asset string, transactionId common.TransactionID) *common.Balance {
	switch transactionId {
	case common.SpotID:
		{
			return b.balanceManager.GetSpotBalance(asset)
		}
	case common.FutureID:
		{
			return b.balanceManager.GetFutureBalance(asset)

		}

	default:
		panic("BinaAccountManager::GetBalance invalid transactionId : " + transactionId)
	}
}

func (b *BinaAccountManager) GetFuturePosition(symbol string, transactionId common.TransactionID) *common.FuturePosition {
	return b.balanceManager.GetFuturePosition(symbol, transactionId)
}

func (b *BinaAccountManager) GetOrders(symbol string, transactionId common.TransactionID) []*common.Order {

	if v, ok := b.orderManager[transactionId].OpenOrdersBySymbol(symbol); ok {
		return v
	} else {
		return nil
	}
}

func (b *BinaAccountManager) RegisterSystemSymbols(symbols []string) {
	RegisterSystemSymbols(symbols)
}

func (b *BinaAccountManager) WsUpdateSpotBalance(balance common.SpotBalance) {
	b.balanceManager.WsUpdateSpotBalance(balance)
}

func (b *BinaAccountManager) WsUpdateOrderOnTrade(info common.OrderTradeUpdateInfo) {
	b.orderManager[info.Transaction].WsUpdateOrderOnTrade(info)
}

func (b *BinaAccountManager) WsUpdateFutureBalancePosition(balancePosition common.WsFutureBalancePosition) {
	balances := balancePosition.FutureBalances
	positions := balancePosition.FuturePositions
	b.balanceManager.WsUpdateFutureBalancePosition(balances, positions)
}

func (b *BinaAccountManager) WsUpdateCoinFutureBalancePosition(balancePosition common.WsFutureBalancePosition) {
	balances := balancePosition.FutureBalances
	positions := balancePosition.FuturePositions
	b.balanceManager.WsUpdateCoinFutureBalancePosition(balances, positions)
}

func (b *BinaAccountManager) WsUpdateOrderOnOrder(info common.OrderTradeUpdateInfo) {
	switch info.Transaction {
	case common.SpotID:
		{
			b.wsUpdateSpotOrderOnOrder(info)
		}
	case common.FutureID:
		{
			b.wsUpdateFutureOrderOnOrder(info)
		}

	default:
		panic("BinaAccountManager::WsUpdateOrderOnOrder unknown transaction " + info.Transaction)
	}
}

func (b *BinaAccountManager) wsUpdateSpotOrderOnOrder(info common.OrderTradeUpdateInfo) {
	b.orderManager[common.SpotID].WsUpdateOrderOnOrder(info)

	if info.Status == common.NEW && (strings.HasPrefix(info.Cid, "electron") || strings.HasPrefix(info.Cid, "web")) {
		b.restSetSpotOrder()
		b.restSpotBalance()
	}
}

func (b *BinaAccountManager) wsUpdateFutureOrderOnOrder(info common.OrderTradeUpdateInfo) {
	b.orderManager[common.FutureID].WsUpdateOrderOnOrder(info)

	if info.Status == common.NEW && (strings.HasPrefix(info.Cid, "electron") || strings.HasPrefix(info.Cid, "web")) {
		b.restSetFutureOrder()
		b.restSetFutureBalancePosition()
	}
}

func (b *BinaAccountManager) initAccountInfo() {
	b.restSetSpotBalance()
	b.restSetFutureBalancePosition()
	b.restSetSpotOrder()
	b.restSetFutureOrder()
}

func (b *BinaAccountManager) restSpotBalance() (cmap.ConcurrentMap, error) {
	if rawBalances, err := GetSpotBalance(b.spotCli); err == nil {
		balances := cmap.New()
		for k, v := range rawBalances {
			balances.Set(k, v)
		}
		return balances, nil
	} else {
		return nil, err
	}
}

func (b *BinaAccountManager) restSetSpotBalance() {
	if balances, err := b.restSpotBalance(); err == nil {
		b.balanceManager.SetSpotBalance(balances)
	}
}

func (b *BinaAccountManager) restSpotOrder() ([]*common.Order, error) {
	openOrders, err := b.spotCli.NewListOpenOrdersService().
		Do(context.Background())
	if err != nil {
		common.Logger.Error("Binance restSpotOrder Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()))
		return nil, err
	}
	orders := make([]*common.Order, 0, len(openOrders))
	for _, o := range openOrders {

		if symbol, err := TryGetSysSymbol(o.Symbol); err == nil {
			oo := common.Order{
				Exchange:    common.BINANCEID,
				Transaction: common.SpotID,
				Symbol:      symbol,
				Id:          common.Int2Str(o.OrderID),
				Cid:         o.ClientOrderID,
				Side:        common.SideID(o.Side),
				Type:        common.OrderTypeID(o.Type),
				FilledSize:  common.Str2Float(o.ExecutedQuantity),
				Size:        common.Str2Float(o.OrigQuantity),
				Price:       common.Str2Float(o.Price),
				CreateTime:  o.Time,
				Status:      common.Normal}
			orders = append(orders, &oo)
		} else {
			common.Logger.Warn("Binance restSpotOrder unknown symbol ",
				zap.String("symbol", o.Symbol))
		}
	}
	return orders, nil
}

func (b *BinaAccountManager) restSetSpotOrder() {
	if orders, err := b.restSpotOrder(); err == nil {
		b.orderManager[common.SpotID].SetOpenOrder(orders)
	}
}

func (b *BinaAccountManager) restFutureBalancePosition() (cmap.ConcurrentMap, cmap.ConcurrentMap, error) {
	res, err := b.futureCli.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, nil, err
	}
	balances := cmap.New()
	for _, r := range res.Assets {
		balances.Set(r.Asset, &common.Balance{
			Balance:          common.Str2Float(r.WalletBalance),
			AvailableBalance: common.Str2Float(r.AvailableBalance),
		})
	}
	positions := cmap.New()
	for _, r := range res.Positions {
		symbol, err := TryGetSysSymbol(r.Symbol)
		if err != nil {
			continue
		}
		pos, ok := positions.Get(symbol)
		if !ok {
			positions.Set(symbol, &common.FuturePosition{})
			pos, _ = positions.Get(symbol)
		}
		switch r.PositionSide {
		case futures.PositionSideTypeLong:
			pos.(*common.FuturePosition).LONG = &common.SidePosition{
				Amount:        common.Str2Float(r.PositionAmt),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedProfit),
			}
		case futures.PositionSideTypeShort:
			pos.(*common.FuturePosition).SHORT = &common.SidePosition{
				Amount:        common.Str2Float(r.PositionAmt),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedProfit),
			}
		case futures.PositionSideTypeBoth:
			pos.(*common.FuturePosition).BOTH = &common.SidePosition{
				Amount:        common.Str2Float(r.PositionAmt),
				EntryPrice:    common.Str2Float(r.EntryPrice),
				UnrealizedPnL: common.Str2Float(r.UnrealizedProfit),
			}
		}
	}
	return balances, positions, nil
}

func (b *BinaAccountManager) restSetFutureBalancePosition() {
	if balances, positions, err := b.restFutureBalancePosition(); err == nil {
		b.balanceManager.SetFutureBalancePosition(balances, positions)
	}
}

func (b *BinaAccountManager) restFutureOrder() ([]*common.Order, error) {
	openOrders, err := b.futureCli.NewListOpenOrdersService().
		Do(context.Background())
	if err != nil {
		common.Logger.Error("BinaAccountManager::restFutureOrder Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()))
		return nil, err
	}
	orders := make([]*common.Order, len(openOrders))
	for i, o := range openOrders {
		oo := common.Order{
			Exchange:     common.BINANCEID,
			Transaction:  common.FutureID,
			Symbol:       GetSysSymbol(o.Symbol),
			Id:           common.Int2Str(o.OrderID),
			Cid:          o.ClientOrderID,
			Side:         common.SideID(o.Side),
			PositionSide: common.PositionID(o.PositionSide),
			Type:         common.OrderTypeID(o.Type),
			FilledSize:   common.Str2Float(o.ExecutedQuantity),
			Size:         common.Str2Float(o.OrigQuantity),
			Price:        common.Str2Float(o.Price),
			CreateTime:   o.Time,
			Status:       common.Normal}
		orders[i] = &oo
	}
	return orders, nil
}

func (b *BinaAccountManager) restSetFutureOrder() {
	if orders, err := b.restFutureOrder(); err == nil {
		b.orderManager[common.FutureID].SetOpenOrder(orders)
	}
}

func (b *BinaAccountManager) createLimitTypeSpotOrderProcess(data *common.Order) (string, error) {
	var (
		res *binance.CreateOrderResponse
		id  string
		err error
	)
	symbol := GetExchangeSymbol(data.Symbol)

	service := b.spotCli.NewCreateOrderService().NewClientOrderID(data.Cid).Symbol(symbol).
		Side(binance.SideType(data.Side)).Type(binance.OrderType(data.Type)).Quantity(common.Float2Str(data.Size)).
		Price(common.Float2Str(data.Price))

	if data.Type == common.OrderTypeLimitMaker {
		res, err = service.Do(context.Background())
	} else if data.Type == common.OrderTypeLimit {
		res, err = service.TimeInForce(binance.TimeInForceTypeGTC).Do(context.Background())
	} else {
		err = errors.New("createLimitTypeSpotOrderProcess: Unknown data.Type")
	}
	common.Logger.Info("Binance createLimitTypeSpotOrderProcess ",
		zap.Int("accountIndex", int(b.accountIndex)),
		zap.Any("order", data))

	if err != nil {
		common.Logger.Error("Binance createLimitTypeSpotOrderProcess Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()),
			zap.Any("data", data))

		b.systemAgent.OnError(common.ErrorMsg{
			Exchange:     common.BINANCEID,
			Transaction:  common.SpotID,
			AccountIndex: b.accountIndex,
			ActionID:     common.CREATE_LIMITTYPE_SPOT_ORDER,
			Symbol:       data.Symbol,
			Id:           data.Id,
			Cid:          data.Cid,
			Side:         data.Side,
			Size:         data.Size,
			Error:        err.Error(),
			Timestamp:    common.SystemNanoSeconds()})

	} else {
		data.CreateTime = res.TransactTime
		data.Cid = res.ClientOrderID
		id = common.Int2Str(res.OrderID)
	}
	return id, err
}

func (b *BinaAccountManager) createLimitTypeFutureOrderProcess(data *common.Order) (string, error) {
	var (
		res *futures.CreateOrderResponse
		id  string
		err error
	)
	symbol := GetExchangeSymbol(data.Symbol)

	orderService := b.futureCli.NewCreateOrderService().NewClientOrderID(data.Cid).Symbol(symbol).
		PositionSide(futures.PositionSideType(data.PositionSide)).Side(futures.SideType(data.Side)).
		Type(futures.OrderType(common.OrderTypeLimit)).Quantity(common.Float2Str(data.Size)).Price(common.Float2Str(data.Price))

	switch data.Type {
	case common.OrderTypeLimitMaker:
		orderService = orderService.TimeInForce(futures.TimeInForceTypeGTX)
	case common.OrderTypeLimit:
		orderService = orderService.TimeInForce(futures.TimeInForceTypeGTC)
	}

	if data.ReduceOnly {
		orderService = orderService.ReduceOnly(data.ReduceOnly)
	}

	res, err = orderService.Do(context.Background())

	common.Logger.Info("Binance createLimitTypeFutureOrderProcess ",
		zap.Int("accountIndex", int(b.accountIndex)),
		zap.Any("order", data))

	if err != nil {
		common.Logger.Error("Binance createLimitTypeFutureOrderProcess Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()),
			zap.Any("data", data))
		b.systemAgent.OnError(common.ErrorMsg{
			Exchange:     common.BINANCEID,
			Transaction:  common.FutureID,
			AccountIndex: b.accountIndex,
			ActionID:     common.CREATE_LIMITTYPE_FUTURE_ORDER,
			Symbol:       data.Symbol,
			Id:           data.Id,
			Cid:          data.Cid,
			Side:         data.Side,
			Size:         data.Size,
			Error:        err.Error(),
			Timestamp:    common.SystemNanoSeconds()})
	} else {
		data.CreateTime = res.UpdateTime
		data.Cid = res.ClientOrderID
		id = common.Int2Str(res.OrderID)
	}
	return id, err
}

func (b *BinaAccountManager) createMarketSpotOrderProcess(data *common.Order) (id string, err error) {
	symbol := GetExchangeSymbol(data.Symbol)
	res, err := b.spotCli.NewCreateOrderService().NewClientOrderID(data.Cid).Symbol(symbol).
		Side(binance.SideType(data.Side)).Type(binance.OrderType(data.Type)).Quantity(common.Float2Str(data.Size)).Do(context.Background())

	common.Logger.Info("Binance createMarketSpotOrderProcess ",
		zap.Int("accountIndex", int(b.accountIndex)),
		zap.Any("order", data))

	if err != nil {
		common.Logger.Error("Binance createMarketSpotOrderProcess Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()),
			zap.Any("data", data))
		b.systemAgent.OnError(common.ErrorMsg{
			Exchange:     common.BINANCEID,
			Transaction:  common.SpotID,
			AccountIndex: b.accountIndex,
			ActionID:     common.CREATE_MARKET_SPOT_ORDER,
			Symbol:       data.Symbol,
			Id:           data.Id,
			Cid:          data.Cid,
			Side:         data.Side,
			Size:         data.Size,
			Error:        err.Error(),
			Timestamp:    common.SystemNanoSeconds()})
	} else {
		data.CreateTime = res.TransactTime
		data.Cid = res.ClientOrderID
		id = common.Int2Str(res.OrderID)
	}
	return id, err
}

func (b *BinaAccountManager) cancelSpotOrderProcess(data common.CancelInfo) error {
	symbol := GetExchangeSymbol(data.Symbol)
	response, err := b.spotCli.NewCancelOrderService().OrderID(common.Str2Int(data.Id)).Symbol(symbol).Do(context.Background())

	common.Logger.Info("Binance cancelSpotOrderProcess ",
		zap.Int("accountIndex", int(b.accountIndex)),
		zap.Any("common.CancelInfo", data))

	if err != nil {
		if o, ok := b.orderManager[common.SpotID].OpenOrder(data.Id); ok {
			b.systemAgent.OnError(common.ErrorMsg{
				Exchange:     common.BINANCEID,
				Transaction:  common.SpotID,
				AccountIndex: b.accountIndex,
				ActionID:     common.CANCEL_SPOT_ORDER,
				Symbol:       data.Symbol,
				Cid:          o.Cid,
				Id:           data.Id,
				Side:         o.Side,
				Size:         o.Size,
				Error:        err.Error(),
				Timestamp:    common.SystemNanoSeconds()})
		} else {
			b.systemAgent.OnError(common.ErrorMsg{
				Exchange:     common.BINANCEID,
				Transaction:  common.SpotID,
				AccountIndex: b.accountIndex,
				ActionID:     common.CANCEL_SPOT_ORDER,
				Symbol:       data.Symbol,
				Id:           data.Id,
				Side:         "",
				Size:         0,
				Error:        err.Error(),
				Timestamp:    common.SystemNanoSeconds()})
		}
		common.Logger.Error("Binance cancelSpotOrderProcess Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()),
			zap.String("cancelOrder", common.Marshal(data)),
			zap.String("response", common.Marshal(response)))
		if !strings.Contains(err.Error(), "Unknown order") {
			return err
		}
	}
	return nil
}

func (b *BinaAccountManager) cancelAllSpotOrderProcess(data common.CancelInfo) {
	symbol := GetExchangeSymbol(data.Symbol)
	_, err := b.spotCli.NewCancelOpenOrdersService().Symbol(symbol).Do(context.Background())

	common.Logger.Info("Binance cancelAllSpotOrderProcess ",
		zap.Int("accountIndex", int(b.accountIndex)),
		zap.Any("common.CancelInfo", data))

	if err != nil {
		common.Logger.Error("Binance cancelAllSpotOrderProcess Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()),
			zap.Any("data", data))
	}
}

func (b *BinaAccountManager) createMarketFutureOrderProcess(data *common.Order) (id string, err error) {
	symbol := GetExchangeSymbol(data.Symbol)

	orderService := b.futureCli.NewCreateOrderService().NewClientOrderID(data.Cid).Symbol(symbol).
		Side(futures.SideType(data.Side)).Type(futures.OrderType(data.Type)).Quantity(common.Float2Str(data.Size))

	if data.ReduceOnly {
		orderService = orderService.ReduceOnly(data.ReduceOnly)
	}

	res, err := orderService.Do(context.Background())

	common.Logger.Info("Binance createMarketFutureOrderProcess ",
		zap.Int("accountIndex", int(b.accountIndex)),
		zap.Any("order", data))

	if err != nil {
		common.Logger.Error("Binance createMarketFutureOrderProcess Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()),
			zap.Any("data", data))
		b.systemAgent.OnError(common.ErrorMsg{
			Exchange:     common.BINANCEID,
			Transaction:  common.FutureID,
			AccountIndex: b.accountIndex,
			ActionID:     common.CREATE_MARKET_FUTURE_ORDER,
			Symbol:       data.Symbol,
			Id:           data.Id,
			Cid:          data.Cid,
			Side:         data.Side,
			Size:         data.Size,
			Error:        err.Error(),
			Timestamp:    common.SystemNanoSeconds()})
	} else {
		data.CreateTime = res.UpdateTime
		data.Cid = res.ClientOrderID
		id = common.Int2Str(res.OrderID)
	}
	return id, err
}

func (b *BinaAccountManager) cancelFutureOrderProcess(data common.CancelInfo) error {
	symbol := GetExchangeSymbol(data.Symbol)
	response, err := b.futureCli.NewCancelOrderService().OrderID(common.Str2Int(data.Id)).Symbol(symbol).Do(context.Background())

	common.Logger.Info("Binance cancelFutureOrderProcess ",
		zap.Int("accountIndex", int(b.accountIndex)),
		zap.Any("common.CancelInfo", data))

	if err != nil {
		if o, ok := b.orderManager[common.FutureID].OpenOrder(data.Id); ok {
			b.systemAgent.OnError(common.ErrorMsg{
				Exchange:     common.BINANCEID,
				Transaction:  common.FutureID,
				AccountIndex: b.accountIndex,
				ActionID:     common.CANCEL_FUTURE_ORDER,
				Symbol:       data.Symbol,
				Id:           data.Id,
				Cid:          o.Cid,
				Side:         o.Side,
				Size:         o.Size,
				Error:        err.Error(),
				Timestamp:    common.SystemNanoSeconds()})
		} else {
			b.systemAgent.OnError(common.ErrorMsg{
				Exchange:     common.BINANCEID,
				Transaction:  common.FutureID,
				AccountIndex: b.accountIndex,
				ActionID:     common.CANCEL_FUTURE_ORDER,
				Symbol:       data.Symbol,
				Id:           data.Id,
				Side:         "",
				Size:         0,
				Error:        err.Error(),
				Timestamp:    common.SystemNanoSeconds()})
		}
		common.Logger.Error("Binance cancelFutureOrderProcess Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()),
			zap.String("cancelOrder", common.Marshal(data)),
			zap.String("response", common.Marshal(response)))
		if !strings.Contains(err.Error(), "Unknown order") {
			return err
		}
	}
	return nil
}

func (b *BinaAccountManager) cancelFutureAllOrderProcess(data common.CancelInfo) {
	symbol := GetExchangeSymbol(data.Symbol)
	err := b.futureCli.NewCancelAllOpenOrdersService().Symbol(symbol).Do(context.Background())

	common.Logger.Info("Binance cancelFutureAllOrderProcess ",
		zap.Int("accountIndex", int(b.accountIndex)),
		zap.Any("common.CancelInfo", data))

	if err != nil {
		common.Logger.Error("BinaAccountManager::cancelFutureAllOrderProcess Error",
			zap.Int("accountIndex", int(b.accountIndex)),
			zap.String("error", err.Error()),
			zap.Any("data", data))
	}
}

func (b *BinaAccountManager) onTimer() {
	// 每60秒进行一次账户和订单数据同步
	ticker1 := time.NewTicker(Rest_update_interval)
	for {
		<-ticker1.C
	}
}
