package main

import (
	"time"

	"github.com/xh23123/IT_hftcommon/pkg/common"
	bina "github.com/xh23123/IT_hftcommon/pkg/exampleExchangeGateway/binance"
)

func InitWebsocketConnections(systemAgent common.TradeSystemAgent) {
	// 初始化币安账户
	accounts := bina.GetBinanceAccountConfigs(systemAgent)
	for _, v := range accounts {
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
	for _, v := range accounts {
		// 初始化行情Ws
		bina.NewMarketWs(systemAgent, v)
	}
}

func NewRestClient(config map[string]string) common.RestClientInterface {
	return bina.NewRestClient(config)
}
