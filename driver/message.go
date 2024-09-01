package driver

import (
	"encoding/binary"
	"go-driver/pb"
	"io"

	"google.golang.org/protobuf/proto"
)

type MethodDesc struct {
	MethodName string
	Request    proto.Message
	Response   proto.Message
}

type MethodDescList []MethodDesc

var MethodDescs = MethodDescList{
	{MethodName: "Ping", Request: &pb.PingRequest{}, Response: &pb.PingResponse{}},
	{MethodName: "Login", Request: &pb.LoginRequest{}, Response: &pb.LoginResponse{}},
	{MethodName: "Chat", Request: &pb.ChatRequest{}, Response: &pb.ChatResponse{}},
}

type Message []byte

// 从数据流中获取协议号
func (x Message) Method() uint16 {
	return binary.BigEndian.Uint16(x[2:])
}

// 从数据流中获取包体
func (x Message) Payload() []byte {
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

func (x *Message) Reset() {
	if cap(*x) <= 0 {
		return
	}
	*x = (*x)[:0]
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
