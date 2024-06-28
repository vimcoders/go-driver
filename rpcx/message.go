package rpcx

import (
	"bufio"
	"encoding/binary"
	"fmt"

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
func encode(seqNumber uint32, kind uint16, message proto.Message) ([]byte, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var header [10]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(header)+len(b)))
	binary.BigEndian.PutUint32(header[4:], seqNumber)
	binary.BigEndian.PutUint16(header[8:], uint16(kind))
	return append(header[:], b...), nil
}

func (x Message) SeqNumber() uint32 {
	return binary.BigEndian.Uint32(x[4:])
}

// 从数据流中获取协议号
func (x Message) Kind() uint16 {
	return binary.BigEndian.Uint16(x[8:])
}

// 从数据流中获取包体
func (x Message) Message() []byte {
	return x[10:]
}
