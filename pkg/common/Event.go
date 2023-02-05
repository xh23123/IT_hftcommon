package common

// 策略返回 统一数据格式
type ActionEvent struct {
	Exchange     ExchangeID    `json:"exid"`
	AccountIndex AccountIdx    `json:"accountIdx"`
	Transaction  TransactionID `json:"tid"`
	Action       ActionID      `json:"aid"`
	Data         interface{}   `json:"data"`
}
