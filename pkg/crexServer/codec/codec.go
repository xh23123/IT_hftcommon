package codec

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"

	"github.com/xh23123/IT_hftcommon/pkg/crexServer/codec/message"
	"github.com/xh23123/IT_hftcommon/pkg/crexServer/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var messagePool = sync.Pool{New: func() any {
	return new(message.ProxyReq)
}}

var encoderPool = sync.Pool{New: func() any {
	return make([]byte, 1024)
}}

type CodecInterface interface {
	Encode(rsp *message.ProxyRsp, isLittleEndian bool) ([]byte, error)
	Decode([]byte) (*message.ProxyReq, error)
}

var DefaultCodec = defaultCodec{}

type defaultCodec struct {
}

func (c *defaultCodec) Encode(ProxyRsp *message.ProxyRsp, isLittleEnd bool) ([]byte, error) {
	returnData, err := proto.Marshal(ProxyRsp)
	if err != nil {
		return nil, err
	}
	l := len(returnData)
	response := make([]byte, 0, l+headerSize)
	if isLittleEnd {
		response = binary.LittleEndian.AppendUint32(response, uint32(l))
	} else {
		response = binary.BigEndian.AppendUint32(response, uint32(l))
	}
	response = append(response, returnData...)
	return response, nil
}

func (c *defaultCodec) Decode(data []byte) (ret *message.ProxyReq, err error) {
	ret = messagePool.Get().(*message.ProxyReq)
	err = proto.Unmarshal(data, ret)
	if err != nil {
		messagePool.Put(ret)
		return nil, err
	}
	return
}

func RecycleMessage(message *message.ProxyReq) {
	messagePool.Put(message)
}

type Conn interface {
	io.Reader
	// Peek returns the next n bytes without advancing the reader. The bytes stop
	// being valid at the next read call. If Peek returns fewer than n bytes, it
	// also returns an error explaining why the read is short. The error is
	// ErrBufferFull if n is larger than b's buffer size.
	//
	// Note that the []byte buf returned by Peek() is not allowed to be passed to a new goroutine,
	// as this []byte will be reused within event-loop.
	// If you have to use buf in a new goroutine, then you need to make a copy of buf and pass this copy
	// to that new goroutine.
	Peek(n int) (buf []byte, err error)

	// Discard skips the next n bytes, returning the number of bytes discarded.
	//
	// If Discard skips fewer than n bytes, it also returns an error.
	// If 0 <= n <= b.Buffered(), Discard is guaranteed to succeed without
	// reading from the underlying io.Reader.
	Discard(n int) (discarded int, err error)

	// InboundBuffered returns the number of bytes that can be read from the current buffer.
	InboundBuffered() (n int)
}

const headerSize = 4

var (
	ErrNotEnoughHeader = errors.New("not enough header bytes")
	ErrInvalidData     = errors.New("invalid data")
)

func TryGetMessage(c Conn, isLittleEnd bool) (message []byte, err error) {
	inboundBuffered := c.InboundBuffered()
	if inboundBuffered < headerSize {
		return nil, ErrNotEnoughHeader
	}
	headerLen, err := c.Peek(headerSize)
	if err != nil {
		return nil, err
	}
	var dataLen uint32
	if isLittleEnd {
		dataLen = binary.LittleEndian.Uint32(headerLen)
	} else {
		dataLen = binary.BigEndian.Uint32(headerLen)
	}
	peekData, err := c.Peek(int(dataLen) + headerSize)
	if err != nil {
		logger.Error("OnTraffic", zap.Error(err), zap.Uint32("dataLen", dataLen))
		return nil, err
	}
	if len(peekData) != int(dataLen)+headerSize {
		return nil, ErrInvalidData
	} else {
		c.Discard(int(dataLen) + headerSize)
	}
	return peekData[headerSize:], nil
}
