package common

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTradeWsNilExtraPreservesJSON(t *testing.T) {
	trade := TradeWs{
		MarketDataHeader: MarketDataHeader{Exchange: BYBITID, DataID: TradeID,
			Symbol: "BTC-USDT", TransactionID: SpotID, ReceiveTimestamp: 123, ExchangeTimestamp: 456},
		TradeID: 789, TransactionTime: 456, Price: 100.5, Size: 2, TradeCount: 1,
	}
	encoded, err := json.Marshal(trade)
	if err != nil {
		t.Fatal(err)
	}
	// Exact historical shape: adding an unused Extra must not change the payload.
	expected := `{"exid":"6","dataid":"12","symbol":"BTC-USDT","transactionid":"1","recv_timestamp":123,"exchange_timestamp":456,"tradeid":789,"tradetime":456,"price":100.5,"size":2,"tradecount":1,"m":false}`
	if string(encoded) != expected {
		t.Fatalf("unexpected JSON: %s", encoded)
	}
	var decoded TradeWs
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(trade, decoded) {
		t.Fatalf("round-trip mismatch: %#v", decoded)
	}
}

func TestTradeWsExtraJSONRoundTrip(t *testing.T) {
	const id = "20f43950-d8dd-5b31-9112-a178eb6023af"
	for _, tc := range []struct {
		name  string
		extra interface{}
	}{
		{"struct", struct {
			TradeID string `json:"trade_id"`
		}{TradeID: id}},
		{"map", map[string]interface{}{"trade_id": id}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			trade := TradeWs{MarketDataHeader: MarketDataHeader{Extra: tc.extra}, TradeCount: 1}
			encoded, err := json.Marshal(trade)
			if err != nil {
				t.Fatal(err)
			}
			var object map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &object); err != nil {
				t.Fatal(err)
			}
			if string(object["extra"]) != `{"trade_id":"`+id+`"}` {
				t.Fatalf("unexpected extra: %s", object["extra"])
			}
			var decoded TradeWs
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatal(err)
			}
			extra, ok := decoded.Extra.(map[string]interface{})
			if !ok || extra["trade_id"] != id || decoded.TradeCount != 1 {
				t.Fatalf("round-trip mismatch: %#v", decoded)
			}
		})
	}
}
