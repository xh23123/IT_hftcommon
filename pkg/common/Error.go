package common

type ErrorMsg struct {
	Exchange     ExchangeID    `json:"exid"`
	Transaction  TransactionID `json:"tid"`
	AccountIndex AccountIdx    `json:"accountidx"`
	ActionID     ActionID      `json:"aid"`
	Symbol       string        `json:"symbol"`
	Id           string        `json:"id"`
	Cid          string        `json:"cid"`
	Side         SideID        `json:"side"`
	Size         float64       `json:"size"`
	Error        string        `json:"error"`
	Timestamp    int64         `json:"timestamp"`
}
