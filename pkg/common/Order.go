package common

type Order struct {
	Symbol       string      `json:"symbol"`
	Id           string      `json:"id"`
	Cid          string      `json:"cid"`
	Side         SideID      `json:"side"`
	IsIsolated   bool        `json:"is_isolated"`
	PositionSide PositionID  `json:"position_side"`
	Type         OrderTypeID `json:"type"`
	FilledSize   float64     `json:"filled_size"`
	Size         float64     `json:"size"`
	Price        float64     `json:"price"`
	CreateTime   int64       `json:"create_time"`
	CancelTime   int64       `json:"cancel_time"`
	Status       StatusID    `json:"status"`
}

type CancelInfo struct {
	Id         string `json:"id"`
	Cid        string `json:"cid"`
	Symbol     string `json:"symbol"`
	CreateTime int64  `json:"create_time"`
}
