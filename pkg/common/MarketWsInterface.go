package common

type MarketWsInterface interface {
	ResetWs(rs []ResetID)
}

const MarketWsAccountIndex = 0
