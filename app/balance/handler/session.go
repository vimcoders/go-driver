package handler

import (
	"context"
	"go-driver/driver"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/tcp"
	"runtime"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Session struct {
	c tcp.Client
	grpc.ServiceDesc
	total uint64
	unix  int64
	pb.UnimplementedParkourServer
	sync.Pool
}

func (x *Session) ServeTCP(ctx context.Context, buf []byte) error {
	return x.Handle(ctx, buf)
}

func (x *Session) ServeKCP(ctx context.Context, buf []byte) error {
	var request driver.Message = buf
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
	response := x.Pool.Get().(*driver.Message)
	response.WriteUint16(uint16(4 + len(b)))
	response.WriteUint16(method)
	response.Write(b)
	if _, err := response.WriteTo(x.c); err != nil {
		return err
	}
	x.Pool.Put(response)
	return nil
}

func (x *Session) ServeQUIC(ctx context.Context, buf []byte) error {
	return x.Handle(ctx, buf)
}

func (x *Session) Handle(ctx context.Context, buf []byte) error {
	var request driver.Message = buf
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
	response := x.Pool.Get().(*driver.Message)
	response.WriteUint16(uint16(4 + len(b)))
	response.WriteUint16(method)
	response.Write(b)
	if _, err := x.c.Write(*response); err != nil {
		return err
	}
	response.Reset()
	x.Pool.Put(response)
	return nil
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
