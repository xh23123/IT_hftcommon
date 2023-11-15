package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	proxy "github.com/xh23123/IT_hftcommon/pkg/crexServer"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/codec/message"
	"google.golang.org/protobuf/proto"
)

var seq atomic.Int32

func genSeq() string {
	d := seq.Add(1)
	return strconv.Itoa(int(d))
}

func NewWebsocketStartMessage(addr string) *message.ProxyReq {
	var proxyReq message.ProxyReq
	proxyReq.Type = message.Reqtype_WEBSOCKET
	proxyReq.Seq = genSeq()
	var websocketStartReq message.WebsocketStartReq
	websocketStartReq.Url = addr
	m, _ := proto.Marshal(&websocketStartReq)
	proxyReq.Message = m

	return &proxyReq
}

func NewWebsocketWriteMessage(content string) *message.ProxyReq {
	var proxyReq message.ProxyReq
	proxyReq.Type = message.Reqtype_WEBSOCKETWRITE
	proxyReq.Seq = genSeq()
	var websocketWriteReq message.WebsocketWriteReq
	websocketWriteReq.MessageType = websocket.TextMessage
	websocketWriteReq.Message = []byte(content)
	m, _ := proto.Marshal(&websocketWriteReq)
	proxyReq.Message = m

	return &proxyReq
}

func OnMessage(conn net.Conn) error {
	header := make([]byte, 4)
	n, err := conn.Read(header)
	if err != nil || n != 4 {
		log.Fatal("header", err)
	}
	dataLen := 0
	if isLittleEnd {
		dataLen = int(binary.LittleEndian.Uint32(header))
	} else {
		dataLen = int(binary.BigEndian.Uint32(header))
	}
	resp := make([]byte, dataLen)
	n, err = conn.Read(resp)
	if dataLen != n {
		log.Fatal("dataLen != len(resp) - 4", dataLen, len(resp))
	}
	if err != nil {
		log.Fatal("read", err)
	}
	var proxyRsp message.ProxyRsp
	err = proto.Unmarshal(resp, &proxyRsp)
	if err != nil {
		log.Fatal("proto unmarshal proxy rsp", err)
	}
	if proxyRsp.GetError() != "" {
		log.Fatal("proxy err:", proxyRsp.GetError())
	}
	switch proxyRsp.Type {
	case message.Reqtype_WEBSOCKET:
		var resp message.WebsocketStartResp
		err := proto.Unmarshal(proxyRsp.Message, &resp)
		if err != nil {
			return err
		}
		log.Println("Reqtype_WEBSOCKET", resp.String())
	case message.Reqtype_WEBSOCKETPUSH:
		var resp message.WebsocketPushMessage
		err := proto.Unmarshal(proxyRsp.Message, &resp)
		if err != nil {
			return err
		}
		log.Println("Reqtype_WEBSOCKETPUSH", resp.String())
	case message.Reqtype_WEBSOCKETWRITE:
		var resp message.WebsocketWriteRsp
		err := proto.Unmarshal(proxyRsp.Message, &resp)
		if err != nil {
			return err
		}
		log.Println("Reqtype_WEBSOCKETWRITE", resp.String())
	}
	return nil
}

func NewProxyConn(port int) (net.Conn, error) {
	return net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second*2)
}

func TestWs(t *testing.T) {
	port := int(rand.Int31n(1000) + 2000)
	s := proxy.NewServer(port)
	go s.Run()
	defer s.Stop()

	time.Sleep(time.Millisecond * 50)
	conn, err := NewProxyConn(port)
	if err != nil {
		t.Fatal("NewProxyConn", err)
	}
	go func() {
		for {
			OnMessage(conn)
		}
	}()
	// go heartbeat(conn)
	startReq := NewWebsocketStartMessage("wss://stream.bybit.com/v5/public/spot")
	send(conn, startReq)
	writeReq := NewWebsocketWriteMessage(`{
		"op": "subscribe",
		"args": [
			"orderbook.1.BTCUSDT",
			"publicTrade.BTCUSDT",
			"tickers.BTCUSDT",
			"orderbook.50.BTCUSDT"
		]
	}`)
	time.Sleep(time.Millisecond * 50)
	send(conn, writeReq)

	time.Sleep(10 * time.Second)
}

type Bybit struct {
	apiKey        string
	secretKey     string
	ProxyURL      *url.URL
	unifiedStatus int
}
type BybitRestResponse struct {
	RetCode    int64           `json:"retCode"`
	RetMsg     string          `json:"retMsg"`
	Result     json.RawMessage `json:"result"`
	RetExtInfo json.RawMessage `json:"retExtInfo"`
	Time       int64           `json:"time"`
}

var (
	xAPIKey            = "X-BAPI-API-KEY"
	xSign              = "X-BAPI-SIGN"
	xTimestamp         = "X-BAPI-TIMESTAMP"
	xSignType          = "X-BAPI-SIGN-TYPE"
	xRecvWindow        = "X-BAPI-RECV-WINDOW"
	xRecvWindowDefault = "30000"
	xSignTypeDefault   = "2"
)

func (b *Bybit) Sign(timestamp int64, recvWindow, params string) string {
	hmac256 := hmac.New(sha256.New, []byte(b.secretKey))
	hmac256.Write([]byte(strconv.FormatInt(timestamp, 10) + b.apiKey + recvWindow + params))
	return hex.EncodeToString(hmac256.Sum(nil))
}

func HttpRequest(request *http.Request, exchange string, port int) ([]byte, error) {
	conn, err := NewProxyConn(port)
	if err != nil {
		return nil, err
	}
	httpReq := new(message.HttpReq)
	httpReq.Url = request.URL.String()
	httpReq.Method = request.Method
	httpReq.Header = new(message.Header)
	for key, value := range request.Header {
		httpReq.Header.HeaderPair = append(httpReq.Header.HeaderPair, &message.HeaderPair{Key: key, Value: value})

	}
	if request.Body != nil {
		defer request.Body.Close()
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		httpReq.Body = body
	}
	d, _ := proto.Marshal(httpReq)
	proxyReq := new(message.ProxyReq)
	proxyReq.Type = message.Reqtype_HTTP
	proxyReq.Seq = genSeq()
	proxyReq.Message = d

	send(conn, proxyReq)
	if err != nil {
		return nil, err
	}
	resp, err := Recv(conn)
	if err != nil {
		return nil, err
	}

	httpRsp := new(message.HttpResp)
	httpRsp.Reset()
	err = proto.Unmarshal(resp.Message, httpRsp)
	if err != nil {
		return nil, err
	}
	httpData := make([]byte, len(httpRsp.GetResp()))
	copy(httpData, httpRsp.GetResp())
	return httpData, nil
}

func (b *Bybit) Post(path string, body []byte, port int) (*BybitRestResponse, error) {
	endpoint, err := url.JoinPath("https://api-testnet.bybit.com/v5", path)
	if err != nil {
		return nil, fmt.Errorf("url join path fail %s", err)
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("new http request fail %s", err)
	}
	timestamp := time.Now().UnixMilli()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(xAPIKey, b.apiKey)
	req.Header.Set(xSign, b.Sign(timestamp, xRecvWindowDefault, string(body)))
	req.Header.Set(xTimestamp, strconv.FormatInt(timestamp, 10))
	req.Header.Set(xSignType, xSignTypeDefault)
	req.Header.Set(xRecvWindow, xRecvWindowDefault)
	retBody, err := HttpRequest(req, "bybit", port)
	if err != nil {
		return nil, fmt.Errorf("nHttpRequest fail %s", err)
	}
	ret := &BybitRestResponse{}
	err = json.Unmarshal(retBody, ret)
	if err != nil {
		return nil, fmt.Errorf("return body unmarshal fail %s %s", err, string(retBody))
	}
	return ret, nil
}

type OrderCreateRequest struct {
	Category         string `json:"category"`
	Symbol           string `json:"symbol"`
	IsLeverage       int64  `json:"isLeverage,omitempty"`
	Side             string `json:"side"`
	OrderType        string `json:"orderType"`
	Qty              string `json:"qty"`
	Price            string `json:"price,omitempty"`
	TriggerDirection int64  `json:"triggerDirection,omitempty"`
	OrderFilter      string `json:"orderFilter,omitempty"`
	TriggerPrice     string `json:"triggerPrice,omitempty"`
	TriggerBy        string `json:"triggerBy,omitempty"`
	OrderIv          string `json:"orderIv,omitempty"`
	TimeInForce      string `json:"timeInForce,omitempty"`
	PositionIdx      int64  `json:"positionIdx,omitempty"`
	OrderLinkId      string `json:"orderLinkId,omitempty"`
	TakeProfit       string `json:"takeProfit,omitempty"`
	StopLoss         string `json:"stopLoss,omitempty"`
	TpTriggerBy      string `json:"tpTriggerBy,omitempty"`
	SlTriggerBy      string `json:"slTriggerBy,omitempty"`
	ReduceOnly       bool   `json:"reduceOnly,omitempty"`
	CloseOnTrigger   bool   `json:"closeOnTrigger,omitempty"`
	SmpType          string `json:"smpType,omitempty"`
	Mmp              bool   `json:"mmp,omitempty"`
	TpslMode         string `json:"tpslMode,omitempty"`
	TpLimitPrice     string `json:"tpLimitPrice,omitempty"`
	SlLimitPrice     string `json:"slLimitPrice,omitempty"`
	TpOrderType      string `json:"tpOrderType,omitempty"`
	SlOrderType      string `json:"slOrderType,omitempty"`
}

func TestPost(t *testing.T) {
	bybit := Bybit{
		apiKey:    "YCZRJVQXPRLZBFEFHS",
		secretKey: "ZQZTOGJXSUENPSKPIZRXPPBBZHJUXRCLCONJ",
	}
	req := &OrderCreateRequest{
		Category:    "linear",
		Symbol:      "ETHUSDT",
		Side:        "Sell",
		OrderType:   "Market",
		Qty:         "0.01",
		Price:       "0",
		PositionIdx: 0,
	}
	body, _ := json.Marshal(req)
	resp, err := bybit.Post("order/create", body, 8090)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(resp)
}
