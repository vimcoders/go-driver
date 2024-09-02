package handler

import (
	"context"
	"go-driver/driver"
	"go-driver/tcp"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Session struct {
	c tcp.Client
	grpc.ServiceDesc
	sync.Pool
}

func (x *Session) ServeTCP(ctx context.Context, stream []byte) error {
	var request driver.Message = stream
	method, payload := request.Method(), request.Payload()
	dec := func(in any) error {
		if err := proto.Unmarshal(payload, in.(proto.Message)); err != nil {
			return err
		}
		return nil
	}
	reply, err := x.Methods[method].Handler(x, ctx, dec, nil)
	if err != nil {
		return err
	}
	b, err := proto.Marshal(reply.(proto.Message))
	if err != nil {
		return err
	}
	request.Reset()
	request.WriteUint16(uint16(4 + len(b)))
	request.WriteUint16(method)
	request.Write(b)
	if _, err := x.c.Write(request); err != nil {
		return err
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}
