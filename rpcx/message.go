package rpcx

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"

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
func encode(seq, ack uint32, kind uint16, message proto.Message) (Message, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var header [14]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(header)+len(b)))
	binary.BigEndian.PutUint32(header[4:], seq)
	binary.BigEndian.PutUint32(header[8:], ack)
	binary.BigEndian.PutUint16(header[12:], uint16(kind))
	return append(header[:], b...), nil
}

func (x Message) seq() uint32 {
	return binary.BigEndian.Uint32(x[4:])
}

func (x Message) ack() uint32 {
	return binary.BigEndian.Uint32(x[8:])
}

// 从数据流中获取协议号
func (x Message) kind() uint16 {
	return binary.BigEndian.Uint16(x[12:])
}

// 从数据流中获取包体
func (x Message) message() []byte {
	return x[14:]
}

type Pusher struct {
	seq uint32
	ack uint32
	net.Conn
	timeout  time.Duration
	messages []proto.Message
}

func (x *Pusher) push(_ context.Context, iMessage proto.Message) error {
	response, err := x.marshal(iMessage)
	if err != nil {
		return err
	}
	if err := x.SetWriteDeadline(time.Now().Add(x.timeout)); err != nil {
		return err
	}
	if _, err := x.Conn.Write(response); err != nil {
		return err
	}
	return nil
}

func (x *Pusher) marshal(message proto.Message) ([]byte, error) {
	messageName := proto.MessageName(message).Name()
	for i := uint16(0); i < uint16(len(x.messages)); i++ {
		if messageName != proto.MessageName(x.messages[i]).Name() {
			continue
		}
		return encode(x.seq, x.ack, i, message)
	}
	return nil, fmt.Errorf("message %s not registered", proto.MessageName(message))
}
