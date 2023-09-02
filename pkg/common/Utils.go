package common

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func Float2Str(number float64) string {
	num := strconv.FormatFloat(number, 'f', 8, 64)
	for strings.HasSuffix(num, "0") {
		num = strings.TrimSuffix(num, "0")
	}
	num = strings.TrimSuffix(num, ".")
	return num
}

func Str2Float(number string) float64 {
	Price, _ := strconv.ParseFloat(number, 64)
	return Price
}

func Int2Str(n int64) string {
	return strconv.FormatInt(n, 10)
}

func Str2Int(n string) int64 {
	int64Num, _ := strconv.ParseInt(n, 10, 64)
	return int64Num
}

func Round(x float64, precision float64) float64 {

	return math.Round(x*precision) / precision
}

func CheckDecimal(num string) int {
	res := strings.Split(num, ".")
	if len(res) > 1 {
		return len(res[1])
	} else {
		return 0
	}
}

func ChangePrec(num float64, lot int) float64 {
	return math.Round(num*math.Pow10(lot)) / math.Pow10(lot)
}

func SystemMilliSeconds() int64 {
	return time.Now().UnixMilli()
}

func SystemNanoSeconds() int64 {
	return time.Now().UnixNano()
}

func SystemMicroSeconds() int64 {
	return time.Now().UnixMicro()
}

func SystemSeconds() int64 {
	return time.Now().Unix()
}

func GenOrderClientId(exchangeID ExchangeID, accountIndex AccountIdx, dataId DataID, sequence int64) string {
	return fmt.Sprintf("%dX%sX%sX%dX%d", SystemSeconds(), dataId, exchangeID, accountIndex, sequence)
}

func Marshal(v any) string {
	if jsonStr, err := json.Marshal(v); err != nil {
		panic("Marshal failed, err : " + err.Error())
	} else {
		return string(jsonStr)
	}
}
