package tcp

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
	return request, nil
}

// 数据流加密
func encode(kind uint16, message proto.Message) (Message, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var header [6]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(header)+len(b)))
	binary.BigEndian.PutUint16(header[4:], uint16(kind))
	return append(header[:], b...), nil
}

// 从数据流中获取协议号
func (x Message) kind() uint16 {
	return binary.BigEndian.Uint16(x[4:])
}

// 从数据流中获取包体
func (x Message) message() []byte {
	return x[6:]
}

type Pusher struct {
	net.Conn
	timeout  time.Duration
	messages []proto.Message
}

func (x *Pusher) Push(_ context.Context, iMessage proto.Message) error {
	messageName := proto.MessageName(iMessage).Name()
	for i := uint16(0); i < uint16(len(x.messages)); i++ {
		if messageName != proto.MessageName(x.messages[i]).Name() {
			continue
		}
		b, err := encode(i, iMessage)
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
	return fmt.Errorf("%s not registered", messageName)
}
