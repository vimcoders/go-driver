package handle

import (
	"bufio"
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/proto"
)

type Request []byte

// 数据流解密
func decode(b *bufio.Reader) (Request, error) {
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
func encode(kind uint16, message proto.Message) ([]byte, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var header [6]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(header)+len(b)))
	binary.BigEndian.PutUint16(header[4:], uint16(kind))
	return append(header[:], b...), nil
}

// 构造一个来自网络层的模拟数据流
func NewRequest(kind uint16) Request {
	var b [6]byte
	var header Request = b[:]
	binary.BigEndian.PutUint32(header[:], uint32(len(header)))
	binary.BigEndian.PutUint16(header[4:], uint16(kind))
	return header
}

// 从数据流中获取协议号
func (x Request) Kind() uint16 {
	return binary.BigEndian.Uint16(x[4:])
}

// 从数据流中获取包体
func (x Request) Message() []byte {
	return x[6:]
}

// 从数据流中获取应该返回的协议号
func (x Request) Reply() uint16 {
	return x.Kind() + 1
}
