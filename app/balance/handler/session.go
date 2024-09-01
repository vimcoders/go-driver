package handler

import (
	"context"
	"encoding/binary"
	"go-driver/grpcx"
	"go-driver/tcp"
	"io"

	"google.golang.org/protobuf/proto"
)

type Message []byte

// 从数据流中获取协议号
func (x Message) method() uint16 {
	return binary.BigEndian.Uint16(x[2:])
}

// 从数据流中获取包体
func (x Message) payload() []byte {
	return x[4:]
}

func (x *Message) Write(p []byte) (int, error) {
	*x = append(*x, p...)
	return len(p), nil
}

func (x *Message) WriteUint32(v uint32) {
	*x = binary.BigEndian.AppendUint32(*x, v)
}

func (x *Message) WriteUint16(v uint16) {
	*x = binary.BigEndian.AppendUint16(*x, v)
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (x Message) WriteTo(w io.Writer) (n int64, err error) {
	if nBytes := len(x); nBytes > 0 {
		m, e := w.Write(x)
		if m > nBytes {
			panic("bytes.Buffer.WriteTo: invalid Write count")
		}
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// Buffer is now empty; reset.
	return n, nil
}

type MethodDesc struct {
	MethodName string
	Request    proto.Message
	Response   proto.Message
}

type Session struct {
	tcp.Client
	rpc        grpcx.Client
	Token      string
	MethodDesc []MethodDesc
}

func (x *Session) ServeTCP(ctx context.Context, stream []byte) error {
	var request Message = stream
	seq := request.method()
	req := x.MethodDesc[seq].Request
	reply := x.MethodDesc[seq].Response
	methodName := x.MethodDesc[seq].MethodName
	if err := proto.Unmarshal(request.payload(), req); err != nil {
		return err
	}
	if err := x.rpc.Invoke(ctx, methodName, req, reply); err != nil {
		return err
	}
	response, err := encode(seq, reply)
	if err != nil {
		return err
	}
	if _, err := x.Write(response); err != nil {
		return err
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}

// 数据流加密
func encode(seq uint16, message proto.Message) ([]byte, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var iMessage Message
	iMessage.WriteUint16(uint16(4 + len(b)))
	iMessage.WriteUint16(seq)
	if _, err := iMessage.Write(b); err != nil {
		return nil, err
	}
	return iMessage, nil
}
