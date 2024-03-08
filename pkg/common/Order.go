package common

import "math/big"

type Order struct {
	Exchange     ExchangeID      `json:"exid"`
	Transaction  TransactionID   `json:"tid"`
	Symbol       SymbolID        `json:"symbol"`
	Id           OrderidID       `json:"id"`
	Cid          ClientOrderidID `json:"cid"`
	Side         SideID          `json:"side"`
	IsIsolated   bool            `json:"is_isolated"`
	PositionSide PositionID      `json:"position_side"`
	Type         OrderTypeID     `json:"type"`
	FilledSize   float64         `json:"filled_size"`
	Size         float64         `json:"size"`
	Price        float64         `json:"price"`
	CreateTime   int64           `json:"create_time"`
	CancelTime   int64           `json:"cancel_time"`
	ReduceOnly   bool            `json:"reduce_only"`
	Status       StatusID        `json:"status"`
	// dex
	DexAmountIn     big.Int `json:"dex_amount_in"`
	DexMinAmountOut big.Int `json:"dex_amount_out"`
	GasPrice        big.Int `json:"gas_price"`
}

type CancelInfo struct {
	Id         OrderidID       `json:"id"`
	Cid        ClientOrderidID `json:"cid"`
	Symbol     SymbolID        `json:"symbol"`
	CreateTime int64           `json:"create_time"`
}
