package common

// 策略返回 统一数据格式
type ActionEvent struct {
	Exchange     ExchangeID    `json:"exid"`
	AccountIndex AccountIdx    `json:"accountidx"`
	Transaction  TransactionID `json:"tid"`
	Action       ActionID      `json:"aid"`
	Data         interface{}   `json:"data"`
}

// GateWay 统一数据格式
type DataEvent struct {
	ExchangeID    ExchangeID
	AccountIndex  AccountIdx
	TransactionID TransactionID
	DataID        DataID
	Data          interface{}
}

func NewDataEvent(exchangeID ExchangeID, accountIndex AccountIdx, transactionID TransactionID, dataID DataID, symbol SymbolID, data interface{}) *DataEvent {
	return &DataEvent{ExchangeID: exchangeID,
		AccountIndex:  accountIndex,
		TransactionID: transactionID,
		DataID:        dataID,
		Data:          data}
}
