package proxy

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/fasthttp"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/codec"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/codec/message"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/logger"
	"go.uber.org/zap"
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

var proxyRspPoll = sync.Pool{
	New: func() any {
		return &message.ProxyRsp{
			Message: new(message.Message),
		}
	},
}

var websocketPushPool = sync.Pool{
	New: func() any {
		return &message.WebsocketPushMessage{}
	},
}

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

func handleHttpReq(ctx context.Context, httpRequest *message.HttpReq) (msgResp *message.HttpRsp, err error) {
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
	err = HttpClient.DoTimeout(request, httpResp, time.Second*2)
	msgResp = new(message.HttpRsp)
	if err != nil {
		msgResp.Error = err.Error()
		return nil, err
	}
	msgResp.Resp = make([]byte, 0, len(httpResp.Body()))
	msgResp.Resp = append(msgResp.Resp, httpResp.Body()...)
	msgResp.StatusCode = uint32(httpResp.StatusCode())

	return
}

func (c *HttpConn) OnMessage(ctx context.Context, req *message.HttpReq) (*message.HttpRsp, error) {
	return handleHttpReq(ctx, req)
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

func (c *WebSocketConn) OnMessage(ctx context.Context, websocketStartReq *message.WebsocketStartReq) (*message.WebsocketStartRsp, error) {
	dial := websocket.Dialer{
		ReadBufferSize:   defaultReadBufferSize,
		WriteBufferSize:  defaultWriteBufferSize,
		HandshakeTimeout: time.Second * 2,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	websocketConn, _, err := dial.Dial(websocketStartReq.GetUrl(), nil)
	if err != nil {
		return nil, err
	}
	c.wsConn = websocketConn
	c.isConnect = true
	var websocketStartResp message.WebsocketStartRsp
	go c.Forward()
	return &websocketStartResp, nil
}

func clear(proxyRsp *message.ProxyRsp) {
	proxyRsp.Code = 0
	proxyRsp.Type = 0
	proxyRsp.Error = ""
	proxyRsp.Seq = 0
}

func (c *WebSocketConn) forwardWsPushMessage(content []byte, messageType uint32) error {
	proxyRsp := proxyRspPoll.Get().(*message.ProxyRsp)
	defer proxyRspPoll.Put(proxyRsp)
	clear(proxyRsp)
	proxyRsp.Type = message.Reqtype_WEBSOCKETPUSH
	messagePush := websocketPushPool.Get().(*message.WebsocketPushMessage)
	defer websocketPushPool.Put(messagePush)
	messagePush.Reset()
	messagePush.MessageType = messageType
	messagePush.Message = content

	proxyRsp.Message.Body = &message.Message_WebsocketPushMessage{WebsocketPushMessage: messagePush}
	rsp, err := c.codec.Encode(proxyRsp, c.isLittleEndian)
	if err != nil {
		logger.Error("Encode", zap.Error(err))
		return err
	}
	err = send(c.conn, rsp)
	if err != nil {
		logger.Error("send", zap.Error(err))
		return err
	}
	return nil
}

func (c *WebSocketConn) Forward() {
	var err error
	var proxyRsp message.ProxyRsp
	proxyRsp.Type = message.Reqtype_WEBSOCKETPUSH
	defer func() {
		c.wsConn.Close()
	}()
	c.wsConn.SetPongHandler(func(appData string) error {
		return c.forwardWsPushMessage([]byte(appData), websocket.PongMessage)
	})

	c.wsConn.SetPingHandler(func(appData string) error {
		return c.forwardWsPushMessage([]byte(appData), websocket.PingMessage)
	})
	for {
		var messageType int
		var msg []byte
		messageType, msg, err = c.wsConn.ReadMessage()
		if err != nil {
			logger.Error("Forward failed", zap.Error(err))
			return
		}
		logger.Info("Forward", zap.String("message", string(msg)), zap.Int("type", messageType))
		err = c.forwardWsPushMessage(msg, uint32(messageType))
		if err != nil {
			logger.Error("Forward Write failed", zap.Error(err))
			return
		}
	}
}

func (c *WebSocketConn) WriteMessage(ctx context.Context, WebsocketWriteReq *message.WebsocketWriteReq) (*message.WebsocketWriteRsp, error) {
	if c.wsConn == nil {
		return nil, errors.New("websocket not ready")
	}
	err := c.wsConn.WriteMessage(int(WebsocketWriteReq.MessageType), WebsocketWriteReq.Message)
	if err != nil {
		return nil, err
	}

	var websocketWriteRsp message.WebsocketWriteRsp
	return &websocketWriteRsp, nil
}

func (c *WebSocketConn) Close() {
	c.wsConn.Close()
}
