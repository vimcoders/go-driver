// TCP，UDP 接入层加密解密
// TODO:: 熔断，限流，降级
package session

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"go-driver/driver"
	"go-driver/pb"

	"google.golang.org/protobuf/proto"
)

type Handler interface {
	Handle(w driver.ResponsePusher, request, reply proto.Message) error
}

type Request []byte

func ReadRequest(b *bufio.Reader) (request Request, err error) {
	headerBytes, err := b.Peek(int(request.Header()))
	if err != nil {
		return nil, err
	}
	header := binary.BigEndian.Uint32(headerBytes)
	if int(header) > b.Size() {
		return nil, fmt.Errorf("header %v too long", header)
	}
	message, err := b.Peek(int(header))
	if err != nil {
		return nil, err
	}
	return message, nil
}

func NewRequest(kind uint16) Request {
	var b [6]byte
	var header Request = b[:]
	binary.BigEndian.PutUint32(header[:], uint32(len(header)))
	binary.BigEndian.PutUint16(header[4:], uint16(kind))
	return header
}

func (x Request) Header() int32 {
	return 4
}

func (x Request) Kind() uint16 {
	return binary.BigEndian.Uint16(x[4:])
}

func (x Request) Message() []byte {
	return x[6:]
}

type Response []byte

func NewResponse(kind uint16, message proto.Message) (Response, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var header [6]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(header)+len(b)))
	binary.BigEndian.PutUint16(header[4:], uint16(kind))
	return append(header[:], b...), nil
}

type Message []proto.Message

func (x Message) Unmarshal(req []byte) (proto.Message, error) {
	if len(req) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	var request Request = req
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

func (x Message) Marshal(response proto.Message) ([]byte, error) {
	for i := uint16(0); i < uint16(len(x)); i++ {
		if proto.MessageName(response) != proto.MessageName(x[i]) {
			continue
		}
		return NewResponse(i, response)
	}
	return nil, fmt.Errorf("message %s not registered", proto.MessageName(response))
}

var Messages = Message{
	&pb.PingRequest{},
	&pb.PingResponse{},
	&pb.LoginRequest{},
	&pb.LoginResponse{},
}
