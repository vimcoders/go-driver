package handler

import (
	"context"
	"go-driver/driver"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/tcp"
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Session struct {
	c tcp.Client
	grpc.ServiceDesc
	total uint64
	unix  int64
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

func (x *Session) ServeKCP(ctx context.Context, stream []byte) error {
	return x.ServeTCP(ctx, stream)
}

func (x *Session) ServeQUIC(ctx context.Context, stream []byte) error {
	return x.ServeTCP(ctx, stream)
}

func (x *Session) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	unix := time.Now().Unix()
	x.total++
	if unix != x.unix {
		log.Debug(x.total, " request/s", " NumGoroutine ", runtime.NumGoroutine())
		x.total = 0
		x.unix = unix
	}
	return &pb.PingResponse{Message: req.Message}, nil
}

func (x *Session) Close() error {
	return nil
}
