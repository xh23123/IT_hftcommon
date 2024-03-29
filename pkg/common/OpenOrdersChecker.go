package common

import (
	"go.uber.org/zap"
)

type orderChecker struct {
	ExchangeId     ExchangeID
	TransactionId  TransactionID
	OrderCheckFunc func() error
}

type OpenOrderChecker struct {
	needToRecheckOrderFunc []orderChecker
	orderCheckFunc         []orderChecker
	lastFirstCheckTime     int64
}

const firstOrderCheckSecond = 54 // check at 54th second per minute
const secondCheckInterval = 1

func (b *OpenOrderChecker) AddOrderCheckFunc(exchangeId ExchangeID, transactionId TransactionID, checkFunc func() error) {
	if Logger == nil {
		InitLogger("golog/common.log", "info")
	}
	b.orderCheckFunc = append(b.orderCheckFunc, orderChecker{
		ExchangeId:     exchangeId,
		TransactionId:  transactionId,
		OrderCheckFunc: checkFunc,
	})

	if b.needToRecheckOrderFunc == nil {
		b.needToRecheckOrderFunc = []orderChecker{}
	}
}

func (b *OpenOrderChecker) CheckOpenOrdersOnTime(curMilliSec int64) {
	curSecond := curMilliSec / 1000
	if b.firstCheckTime(curSecond) {
		for _, v := range b.orderCheckFunc {
			if err := v.OrderCheckFunc(); err != nil {
				Logger.Error("OpenOrderChecker::CheckOpenOrdersOnTime failed firstCheckTime ",
					zap.String("ExchangeId", string(v.ExchangeId)),
					zap.String("TransactionId", string(v.TransactionId)),
					zap.Error(err))
				b.needToRecheckOrderFunc = append(b.needToRecheckOrderFunc, v)
			} else {
				Logger.Info("OpenOrderChecker::CheckOpenOrdersOnTime success firstCheckTime ",
					zap.String("ExchangeId", string(v.ExchangeId)),
					zap.String("TransactionId", string(v.TransactionId)))
			}
		}
	} else if b.secondCheckTime(curSecond) {
		for _, v := range b.needToRecheckOrderFunc {
			if err := v.OrderCheckFunc(); err != nil {
				Logger.Error("OpenOrderChecker::CheckOpenOrdersOnTime failed secondCheckTime ",
					zap.String("ExchangeId", string(v.ExchangeId)),
					zap.String("TransactionId", string(v.TransactionId)),
					zap.Error(err))
			} else {
				Logger.Info("OpenOrderChecker::CheckOpenOrdersOnTime success secondCheckTime ",
					zap.String("ExchangeId", string(v.ExchangeId)),
					zap.String("TransactionId", string(v.TransactionId)))
			}
		}
		b.needToRecheckOrderFunc = b.needToRecheckOrderFunc[:0]
	}
}

func (b *OpenOrderChecker) firstCheckTime(curSecond int64) bool {

	if (curSecond%60) == firstOrderCheckSecond && b.lastFirstCheckTime != curSecond {
		b.lastFirstCheckTime = curSecond
		return true
	}

	return false
}

func (b *OpenOrderChecker) secondCheckTime(curSecond int64) bool {
	return b.lastFirstCheckTime+secondCheckInterval <= curSecond
}
