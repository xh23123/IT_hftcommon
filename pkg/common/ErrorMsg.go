package common

type ErrorMsg struct {
	Exchange     ExchangeID      `json:"exid"`
	Transaction  TransactionID   `json:"tid"`
	AccountIndex AccountIdx      `json:"accountidx"`
	ActionID     ActionID        `json:"aid"`
	Symbol       SymbolID        `json:"symbol"`
	Id           OrderidID       `json:"id"`
	Cid          OrderClientidID `json:"cid"`
	Side         SideID          `json:"side"`
	Size         float64         `json:"size"`
	Error        error           `json:"error"`
	Timestamp    int64           `json:"timestamp"`
}
