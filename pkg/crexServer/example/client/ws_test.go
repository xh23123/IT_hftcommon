package client

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
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
