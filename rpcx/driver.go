package rpcx

import (
	"bufio"
	"encoding/binary"
	"fmt"

	"go-driver/driver"
	"go-driver/pb"

	"google.golang.org/protobuf/proto"
)

type ResponsePusher = driver.ResponsePusher

type Handler interface {
	ServeRPCX(ResponsePusher, []byte, Option) error
}

type Option []*pb.Option

func (x *Option) Push(opt *pb.Option) {
	this := *x
	for i := 0; i < len(this); i++ {
		if this[i].Key == opt.Key {
			this[i] = opt
			return
		}
	}
	*x = append(*x, opt)
}

func (x Option) Get(key string) string {
	for i := 0; i < len(x); i++ {
		if x[i].Key == key {
			return x[i].Value
		}
	}
	return ""
}

func decode(b *bufio.Reader) (*pb.Message, error) {
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
	var message pb.Message
	if err := proto.Unmarshal(request[4:], &message); err != nil {
		return nil, err
	}
	return &message, nil
}

func encode(message proto.Message) ([]byte, error) {
	request, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(len(b)+len(request)))
	return append(b[:], request...), err
}
