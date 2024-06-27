package rpcx

import (
	"bufio"
	"encoding/binary"
	"errors"
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
func encode(packageNumber uint32, kind uint16, message proto.Message) ([]byte, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var header [10]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(header)+len(b)))
	binary.BigEndian.PutUint32(header[4:], packageNumber)
	binary.BigEndian.PutUint16(header[8:], uint16(kind))
	return append(header[:], b...), nil
}

func (x Message) PackageNumber() uint32 {
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

type ProtoBuf []proto.Message

// 将一个来自底层的二进制流反序列化成一个对象
func (x ProtoBuf) Unmarshal(req []byte) (proto.Message, error) {
	if len(req) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	var request Message = req
	kind := request.Kind()
	if kind >= uint16(len(x)) {
		return nil, fmt.Errorf("message id %v not registered", kind)
	}
	message := x[kind].ProtoReflect().New().Interface()
	if err := proto.Unmarshal(request.Message(), message); err != nil {
		return nil, err
	}
	return message, nil
}
