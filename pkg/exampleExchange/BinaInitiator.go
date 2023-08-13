package main

import (
	"time"

	"github.com/xh23123/IT_hftcommon/pkg/common"
	bina "github.com/xh23123/IT_hftcommon/pkg/exampleExchange/binance"
)

func InitWebsocketConnections(systemAgent common.TradeSystemAgent) {
	// 初始化币安账户
	accounts := bina.GetBinanceAccountConfigs(systemAgent)
	//TODO
	//svc.WsManager[common.BINANCEID] = make([]*common.WsInfo, 0, len(accounts))
	for _, v := range accounts {
		//TODO
		//svc.WsManager[common.BINANCEID] = append(svc.WsManager[common.BINANCEID], &common.WsInfo{})
		// 初始化账户Ws

		bina.NewAccountWs(systemAgent, v)

		time.Sleep(1 * time.Second)
		// 初始化行情Ws
		bina.NewMarketWs(systemAgent, v)
	}
}

func InitMarketWs(systemAgent common.TradeSystemAgent) {
	// 初始化币安账户
	accounts := bina.GetBinanceAccountConfigs(systemAgent)
	//svc.WsManager[common.BINANCEID] = make([]*common.WsInfo, 0, len(accounts))
	for _, v := range accounts {
		//svc.WsManager[common.BINANCEID] = append(svc.WsManager[common.BINANCEID], &common.WsInfo{})
		// 初始化行情Ws
		bina.NewMarketWs(systemAgent, v)
	}
}

func NewRestClient(config map[string]string) common.RestClientInterface {
	return bina.NewRestClient(config)
}
