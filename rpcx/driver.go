package rpcx

import (
	"context"
	"encoding/binary"
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
	Option Option
	net.Conn
}

func (x *Response) Push(ctx context.Context, message proto.Message) (int, error) {
	var opt Option = x.Option
	b, err := proto.Marshal(message)
	if err != nil {
		return 0, err
	}
	response := &pb.Message{
		Message: b,
	}
	response.Option = append(response.Option, &pb.Option{Key: MESSAGEID, Value: opt.Get(MESSAGEID)})
	iResponse, err := proto.Marshal(response)
	if err != nil {
		return 0, err
	}
	buffer := make(driver.Buffer, 4)
	binary.BigEndian.PutUint32(buffer, uint32(len(iResponse)))
	buffer.Write(iResponse)
	if err := x.SetWriteDeadline(time.Now().Add(TIMEOUT)); err != nil {
		return 0, err
	}
	if _, err := x.Conn.Write(buffer); err != nil {
		return 0, err
	}
	return len(iResponse), nil
}
