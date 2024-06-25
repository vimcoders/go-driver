package rpcx

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"

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

type Response struct {
	Option
	net.Conn
	Timeout time.Duration
}

func (x *Response) Push(ctx context.Context, iMessage proto.Message) (int, error) {
	b, err := proto.Marshal(iMessage)
	if err != nil {
		return 0, err
	}
	message := &pb.Message{
		Message: b,
	}
	message.Option = append(message.Option, &pb.Option{Key: MESSAGEID, Value: x.Get(MESSAGEID)})
	response, err := encodeRequest(message)
	if err != nil {
		return 0, err
	}
	if err := x.SetWriteDeadline(time.Now().Add(x.Timeout)); err != nil {
		return 0, err
	}
	if _, err := x.Conn.Write(response); err != nil {
		return 0, err
	}
	return len(response), nil
}

func decodeRequest(b *bufio.Reader) (*pb.Message, error) {
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

func encodeRequest(message proto.Message) ([]byte, error) {
	request, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(len(b)+len(request)))
	return append(b[:], request...), err
}
