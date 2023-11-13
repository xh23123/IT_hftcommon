package proxy

import (
	"context"
	"fmt"
	"time"

	"github.com/panjf2000/gnet/v2"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/codec"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/codec/message"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/logger"
	"go.uber.org/zap"
)

type connKey string

var _ gnet.EventHandler = (*ServerHandler)(nil)

type ServerHandler struct {
	engine      gnet.Engine
	isLittleEnd bool
	codec       codec.CodecInterface
}

func (h *ServerHandler) OnBoot(eng gnet.Engine) (action gnet.Action) {
	h.engine = eng
	return
}
func (h *ServerHandler) OnShutdown(eng gnet.Engine) {
}

func (h *ServerHandler) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	return
}

func (h *ServerHandler) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	return
}

func (h *ServerHandler) response(c gnet.Conn, proxyRsp *message.ProxyRsp) (action gnet.Action) {
	response, err := h.codec.Encode(proxyRsp, h.isLittleEnd)
	if err != nil {
		logger.Error("onTraffic Marshal", zap.Error(err), zap.String("proxyRsp", proxyRsp.String()))
		return gnet.Close
	}
	c.Write(response)
	return gnet.None
}

func (h *ServerHandler) OnTraffic(c gnet.Conn) (action gnet.Action) {
	peekData, err := codec.TryGetMessage(c, h.isLittleEnd)
	logger.Info("OnTraffic", zap.String("message", string(peekData)))
	if err != nil {
		logger.Error("onTraffic get message failed", zap.Error(err))
		return gnet.Close
	}
	msg, err := h.codec.Decode(peekData)
	if err != nil {
		logger.Error("OnTraffic", zap.Error(err))
		return gnet.Close
	}

	defer codec.RecycleMessage(msg)
	switch msg.Type {
	case message.Reqtype_HTTP:
		var conn *HttpConn
		ctx := c.Context()
		if ctx != nil {
			c := ctx.(context.Context)
			conn = c.Value(connKey("conn")).(*HttpConn)
		} else {
			conn = NewHttpConn(c, h.codec)
			ctx := context.Background()
			ctx = context.WithValue(ctx, connKey("conn"), conn)
			c.SetContext(ctx)
		}
		resp, err := conn.OnMessage(context.Background(), msg.GetMessage())
		var proxyRsp message.ProxyRsp
		proxyRsp.Seq = msg.GetSeq()
		proxyRsp.Type = message.Reqtype_HTTP
		proxyRsp.Message = resp
		if err != nil {
			proxyRsp.Error = err.Error()
		}
		return h.response(c, &proxyRsp)

	case message.Reqtype_WEBSOCKET:
		var conn *WebSocketConn
		ctx := c.Context()
		if ctx != nil {
			c := ctx.(context.Context)
			conn = c.Value(connKey("conn")).(*WebSocketConn)
		} else {
			conn = NewWebSocketConn(c, h.codec, h.isLittleEnd)
			ctx := context.Background()
			ctx = context.WithValue(ctx, connKey("conn"), conn)
			c.SetContext(ctx)
		}
		resp, err := conn.OnMessage(context.Background(), msg.GetMessage())
		var proxyRsp message.ProxyRsp
		proxyRsp.Seq = msg.GetSeq()
		proxyRsp.Type = message.Reqtype_WEBSOCKET
		proxyRsp.Message = resp
		if err != nil {
			proxyRsp.Error = err.Error()
		}
		return h.response(c, &proxyRsp)

	case message.Reqtype_WEBSOCKETWRITE:
		ctx := c.Context()
		if ctx == nil {
			logger.Error("onTraffic write too early", zap.String("addr", c.RemoteAddr().String()))
			var proxyRsp message.ProxyRsp
			proxyRsp.Message = []byte("write too early")
			h.response(c, &proxyRsp)
			return
		}
		ctt := ctx.(context.Context)
		conn := ctt.Value(connKey("conn")).(*WebSocketConn)
		if conn == nil {
			logger.Error("onTraffic no conn", zap.String("addr", c.RemoteAddr().String()))
			var proxyRsp message.ProxyRsp
			proxyRsp.Message = []byte("no conn")
			h.response(c, &proxyRsp)
			return
		}
		resp, err := conn.WriteMessage(context.Background(), msg.GetMessage())
		var proxyRsp message.ProxyRsp
		proxyRsp.Seq = msg.GetSeq()
		proxyRsp.Type = message.Reqtype_WEBSOCKETWRITE
		proxyRsp.Message = resp
		if err != nil {
			proxyRsp.Error = err.Error()
		}
		return h.response(c, &proxyRsp)
	case message.Reqtype_HEARTBEAT:
		return h.response(c, &message.ProxyRsp{
			Type:    msg.GetType(),
			Seq:     msg.GetSeq(),
			Message: msg.GetMessage(),
		})
	}
	return
}
func (h *ServerHandler) OnTick() (delay time.Duration, action gnet.Action) {

	return
}

type Server struct {
	PortAddr string
}

func NewServer(port int) *Server {
	return &Server{
		PortAddr: fmt.Sprintf(":%d", port),
	}
}

func (s *Server) Run() error {
	return gnet.Run(&ServerHandler{isLittleEnd: true, codec: &codec.DefaultCodec}, s.PortAddr, gnet.WithTicker(true), gnet.WithLockOSThread(true),
		gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*1), gnet.WithTCPNoDelay(gnet.TCPNoDelay), gnet.WithTCPKeepAlive(time.Second*5))
}

func (s *Server) Stop() error {
	return gnet.Stop(context.Background(), s.PortAddr)
}