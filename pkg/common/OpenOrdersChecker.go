package common

type OpenOrderChecker struct {
	needToRecheckOrderFunc []func() error
	orderCheckFunc         []func() error
	lastFirstCheckTime     int64
}

const firstOrderCheckSecond = 54 // check at 54th second per minute
const secondCheckInterval = 1

func (b *OpenOrderChecker) AddOrderCheckFunc(checkFunc func() error) {
	b.orderCheckFunc = append(b.orderCheckFunc, checkFunc)

	if b.needToRecheckOrderFunc == nil {
		b.needToRecheckOrderFunc = []func() error{}
	}
}

func (b *OpenOrderChecker) CheckOpenOrdersOnTime(curMilliSec int64) {
	curSecond := curMilliSec / 1000
	if b.firstCheckTime(curSecond) {
		for _, v := range b.orderCheckFunc {
			if err := v(); err != nil {
				b.needToRecheckOrderFunc = append(b.needToRecheckOrderFunc, v)
			}
		}
	} else if b.secondCheckTime(curSecond) {
		for _, v := range b.needToRecheckOrderFunc {
			v()
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
