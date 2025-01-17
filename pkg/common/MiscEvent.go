package common

type MiscEvent struct {
	Type  MiscEventTypeID `json:"type"`
	Event interface{}     `json:"event"`
}

type MiscEventTypeID string

const MdConnectionStatus MiscEventTypeID = "MdConnectionStatus"

type MdConnectionStatusEvent struct {
	Exchange      ExchangeID    `json:"exid"`
	DataID        DataID        `json:"dataid"`
	TransactionID TransactionID `json:"transactionid"`
	Connected     bool          `json:"connected"`
	Symbols       []SymbolID    `json:"symbols"`
}
