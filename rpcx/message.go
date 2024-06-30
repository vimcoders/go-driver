package rpcx

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Message []byte

// 数据流解密
func decode(b *bufio.Reader) (Message, error) {
	headerBytes, err := b.Peek(4)
	if err != nil {
		return nil, err
	}
	header := binary.BigEndian.Uint32(headerBytes)
	if int(header) > b.Size() {
		return nil, fmt.Errorf("header %v too long", header)
	}
	request, err := b.Peek(int(header))
	if err != nil {
		return nil, err
	}
	if _, err := b.Discard(len(request)); err != nil {
		return nil, err
	}
	return request, nil
}

// 数据流加密
func encode(seq, ack uint32, kind uint16, iMessage proto.Message) (Message, error) {
	b, err := proto.Marshal(iMessage)
	if err != nil {
		return nil, err
	}
	var message [14]byte
	binary.BigEndian.PutUint32(message[:], uint32(len(message)+len(b)))
	binary.BigEndian.PutUint32(message[4:], seq)
	binary.BigEndian.PutUint32(message[8:], ack)
	binary.BigEndian.PutUint16(message[12:], uint16(kind))
	return append(message[:], b...), nil
}

// 获取请求中的序列号
func (x Message) seq() uint32 {
	return binary.BigEndian.Uint32(x[4:])
}

// 获取请求中的回复确认号
func (x Message) ack() uint32 {
	return binary.BigEndian.Uint32(x[8:])
}

// 从数据流中获取协议号
func (x Message) kind() uint16 {
	return binary.BigEndian.Uint16(x[12:])
}

// 从数据流中获取包体
func (x Message) payload() []byte {
	return x[14:]
}

// 推送一个proto
type Pusher struct {
	seq uint32
	ack uint32
	net.Conn
	timeout    time.Duration
	methodName string
	desc       grpc.ServiceDesc
}

func (x *Pusher) Push(_ context.Context, iMessage proto.Message) error {
	for i := uint16(0); i < uint16(len(x.desc.Methods)); i++ {
		if x.methodName != x.desc.Methods[i].MethodName {
			continue
		}
		b, err := encode(x.seq, x.ack, i, iMessage)
		if err != nil {
			return err
		}
		if err := x.SetWriteDeadline(time.Now().Add(x.timeout)); err != nil {
			return err
		}
		if _, err := x.Conn.Write(b); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("%s not registed", x.methodName)
}
