package client

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	proxy "github.com/xh23123/IT_hftcommon/pkg/crexServer"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/codec/message"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var isLittleEnd = true

func send(c net.Conn, proxyReq *message.ProxyReq) {

	returnData, err := proto.Marshal(proxyReq)
	if err != nil {
		logger.Error("onTraffic Marshal", zap.Error(err), zap.String("proxyRsp", proxyReq.String()))
	}
	l := len(returnData)
	response := make([]byte, 0, l+4)
	if isLittleEnd {
		response = binary.LittleEndian.AppendUint32(response, uint32(l))
	} else {
		response = binary.BigEndian.AppendUint32(response, uint32(l))
	}
	response = append(response, returnData...)
	left := len(response)
	total := len(response)
	for left > 0 {
		n, err := c.Write(response[total-left:])
		if err != nil {
			log.Fatal(err)
		}
		left = left - n
	}
}

func TestHttp(t *testing.T) {
	port := rand.Intn(1000) + 2000
	s := proxy.NewServer(port)
	go s.Run()
	defer s.Stop()
	time.Sleep(time.Microsecond * 10)
	conn, err := NewProxyConn(port)
	if err != nil {
		log.Fatal(err)
	}
	var httpReq message.HttpReq
	httpReq.Url = "https://www.baidu.com"
	proxyReq := message.ProxyReq{
		Message: new(message.Message),
	}
	proxyReq.Type = message.Reqtype_HTTP
	proxyReq.Seq = 1
	proxyReq.Message.Body = &message.Message_HttpReq{HttpReq: &httpReq}

	now := time.Now()
	defer func() {
		log.Println(time.Since(now))
	}()
	send(conn, &proxyReq)
	header := make([]byte, 4)
	n, err := conn.Read(header)
	if err != nil || n != 4 {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	var proxyRsp message.ProxyRsp
	err = proto.Unmarshal(resp, &proxyRsp)
	if err != nil {
		log.Fatal(err)
	}
	httpRsp := proxyRsp.Message.GetHttpRsp()
	log.Println(httpRsp.String())
}

// CONTENTTYPE http content-type
const CONTENTTYPE = "application/json"

const FutureHost = "https://api.coinex.com/perpetual"

// USERAGENT http User-Agent
const USERAGENT = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"

func generateFutureAuthorization(str string) string {
	hash := sha256.Sum256([]byte(str))
	sha := hex.EncodeToString(hash[:])
	return sha
}

func interfaceToString(v interface{}) string {
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Bool:
		return strings.Title(fmt.Sprintf("%v", v))
	case reflect.Slice, reflect.Array:
		var items []string
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			items = append(items, "'"+interfaceToString(s.Index(i))+"'")
		}
		return "[" + strings.Join(items, ", ") + "]"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func FixTimestamp() uint64 {
	return uint64(time.Now().UnixMilli())
}

func httpRequest(method, urlHost string, reqParameters map[string]interface{}, accessId, secretKey string, needSign bool) (*message.ProxyReq, error) {
	url := FutureHost + urlHost
	params := make(map[string]interface{}, len(reqParameters))
	for k, v := range reqParameters {
		params[k] = v
	}
	params["timestamp"] = FixTimestamp()
	queryParamsString := ""
	for i, k := range params {
		queryParamsString += fmt.Sprintf("%s=%s&", i, interfaceToString(k))
	}

	queryParamsString = strings.TrimRight(queryParamsString, "&")
	var reqBody []byte

	var headPairs []*message.HeaderPair
	if needSign {
		toEncodeparamsString := queryParamsString + "&secret_key=" + secretKey
		headPairs = append(headPairs, &message.HeaderPair{Key: "Authorization", Value: []string{generateFutureAuthorization(toEncodeparamsString)}})
	} else {
		for key, value := range params {
			queryParamsString += fmt.Sprintf("%s=%s&", key, interfaceToString(value))
		}
		queryParamsString = strings.TrimRight(queryParamsString, "&")
	}

	proxyReq := message.ProxyReq{
		Type:    message.Reqtype_HTTP,
		Seq:     genSeq(),
		Message: new(message.Message),
	}
	headPairs = append(headPairs, &message.HeaderPair{Key: "Content-Type", Value: []string{CONTENTTYPE}})
	headPairs = append(headPairs, &message.HeaderPair{Key: "User-Agent", Value: []string{USERAGENT}})
	headPairs = append(headPairs, &message.HeaderPair{Key: "AccessId", Value: []string{accessId}})
	if method == "POST" {
		reqBody = []byte(queryParamsString)
	}
	if method == "GET" || method == "DELETE" {
		url = url + "?" + queryParamsString
	}
	httpRes := message.HttpReq{
		Url:     url,
		Body:    reqBody,
		Method:  method,
		Timeout: 2,
		Header: &message.Header{
			HeaderPair: headPairs,
		},
	}
	proxyReq.Message.Body = &message.Message_HttpReq{HttpReq: &httpRes}
	return &proxyReq, nil
}

func Recv(conn net.Conn) (*message.ProxyRsp, error) {
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
	return &proxyRsp, nil
}

func TestOrder(t *testing.T) {
	port := 8094
	s := proxy.NewServer(port)
	go s.Run()
	defer s.Stop()

	params := make(map[string]interface{})
	params["market"] = "DOGEUSDT"
	params["side"] = 2
	params["amount"] = "100"
	params["price"] = "0.065"
	proxyReq, err := httpRequest("POST", "/v1/order/put_limit", params, os.Getenv("COINEX_API_KEY"), os.Getenv("COINEX_SECRET_KEY"), true)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := NewProxyConn(port)
	if err != nil {
		t.Fatal(err)
	}
	send(conn, proxyReq)
	rsp, err := Recv(conn)
	if err != nil {
		t.Fatal(err)
	}
	httpResp := rsp.Message.GetHttpRsp()
	t.Log(httpResp.String())
}
