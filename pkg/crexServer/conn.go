package proxy

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/fasthttp"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/codec"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/codec/message"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type noCopy struct{}

func (n *noCopy) Lock()   {}
func (n *noCopy) Unlock() {}

type Conn interface {
	OnMessage(ctx context.Context, conn gnet.Conn, data []byte) error
	Shutdown() error
}

const (
	defaultReadBufferSize  = 4096
	defaultWriteBufferSize = 4096
)

var HttpClient = &fasthttp.Client{
	NoDefaultUserAgentHeader: true, // Don't send: User-Agent: fasthttp
	MaxConnsPerHost:          100,
	ReadBufferSize:           defaultReadBufferSize,  // Make sure to set this big enough that your whole request can be read at once.
	WriteBufferSize:          defaultWriteBufferSize, // Same but for your response.
	ReadTimeout:              8 * time.Second,
	WriteTimeout:             8 * time.Second,
	MaxIdleConnDuration:      30 * time.Second,

	DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this.
}

type HttpConn struct {
	noCopy
	c          gnet.Conn
	codec      codec.CodecInterface
	httpClient *fasthttp.Client
}

func NewHttpConn(c gnet.Conn, codec codec.CodecInterface) *HttpConn {
	return &HttpConn{
		c:          c,
		httpClient: &fasthttp.Client{},
		codec:      codec,
	}
}

func handleHttpReq(ctx context.Context, m []byte) (resp []byte, err error) {
	httpRequest := new(message.HttpReq)
	err = proto.Unmarshal(m, httpRequest)
	if err != nil {
		return nil, err
	}

	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	for _, headerPair := range httpRequest.GetHeader().GetHeaderPair() {
		for _, value := range headerPair.Value {
			request.Header.Add(headerPair.Key, value)
		}
	}
	request.AppendBody([]byte(httpRequest.GetBody()))
	request.SetRequestURI(httpRequest.GetUrl())
	request.Header.SetMethod(httpRequest.Method)

	if err != nil {
		return nil, err
	}
	httpResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(httpResp)
	err = HttpClient.DoTimeout(request, httpResp, time.Second*HttpClient.ReadTimeout)
	msgResp := new(message.HttpResp)
	if err != nil {
		msgResp.Error = err.Error()
		return nil, err
	}
	msgResp.Resp = make([]byte, 0, len(httpResp.Body()))
	msgResp.Resp = append(msgResp.Resp, httpResp.Body()...)
	msgResp.StatusCode = uint32(httpResp.StatusCode())

	return proto.Marshal(msgResp)
}

func (c *HttpConn) OnMessage(ctx context.Context, data []byte) ([]byte, error) {
	return handleHttpReq(ctx, data)
}

func (c *HttpConn) Shutdown() error {

	return nil
}

type WebSocketConn struct {
	noCopy
	conn           gnet.Conn
	wsConn         *websocket.Conn
	isConnect      bool
	codec          codec.CodecInterface
	isLittleEndian bool
}

func NewWebSocketConn(conn gnet.Conn, codec codec.CodecInterface, isLittleEndian bool) *WebSocketConn {
	return &WebSocketConn{
		conn:           conn,
		codec:          codec,
		isLittleEndian: isLittleEndian,
	}
}

func (c *WebSocketConn) OnMessage(ctx context.Context, data []byte) ([]byte, error) {
	var WebsocketStartReq message.WebsocketStartReq
	err := proto.Unmarshal(data, &WebsocketStartReq)
	if err != nil {
		return nil, err
	}
	dial := websocket.Dialer{
		ReadBufferSize:   defaultReadBufferSize,
		WriteBufferSize:  defaultWriteBufferSize,
		HandshakeTimeout: time.Second * 2,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	websocketConn, _, err := dial.Dial(WebsocketStartReq.GetUrl(), nil)
	if err != nil {
		return nil, err
	}
	c.wsConn = websocketConn
	c.isConnect = true
	var websocketStartResp message.WebsocketStartResp
	go c.Forward()
	return proto.Marshal(&websocketStartResp)
}

func (c *WebSocketConn) Forward() {
	defer func() {
		c.conn.Close()
		c.wsConn.Close()
	}()
	for {
		messageType, msg, err := c.wsConn.ReadMessage()
		if err != nil {
			logger.Error("Forward failed", zap.Error(err))
			return
		}
		logger.Info("Forward", zap.String("message", string(msg)), zap.Int("type", messageType))
		var messagePush message.WebsocketPushMessage
		messagePush.Message = msg
		messagePush.MessageType = uint32(messageType)
		data, err := proto.Marshal(&messagePush)
		if err != nil {
			logger.Error("Forward failed", zap.Error(err))
			return
		}
		var proxyRsp message.ProxyRsp
		proxyRsp.Message = data
		proxyRsp.Type = message.Reqtype_WEBSOCKETPUSH
		rsp, err := c.codec.Encode(&proxyRsp, c.isLittleEndian)
		if err != nil {
			logger.Error("Forward failed", zap.Error(err))
			return
		}
		c.conn.Write(rsp)
		if err != nil {
			logger.Error("Forward Write failed", zap.Error(err))
			return
		}
	}
}

func (c *WebSocketConn) WriteMessage(ctx context.Context, data []byte) ([]byte, error) {
	var WebsocketWriteReq message.WebsocketWriteReq
	err := proto.Unmarshal(data, &WebsocketWriteReq)
	if err != nil {
		return nil, err
	}
	err = c.wsConn.WriteMessage(int(WebsocketWriteReq.MessageType), WebsocketWriteReq.Message)
	if err != nil {
		return nil, err
	}

	var websocketWriteRsp message.WebsocketWriteRsp
	return proto.Marshal(&websocketWriteRsp)
}
