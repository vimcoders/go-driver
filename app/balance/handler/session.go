package handler

import (
	"context"
	"go-driver/log"
	"go-driver/pb"
	"runtime"
	"time"
)

type Session struct {
	total uint64
	unix  int64
	pb.UnimplementedParkourServer
}

func (x *Session) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	unix := time.Now().Unix()
	x.total++
	if unix != x.unix {
		log.Debug(x.total, " request/s", " NumGoroutine ", runtime.NumGoroutine(), req.Message)
		x.total = 0
		x.unix = unix
	}
	return &pb.PingResponse{Message: req.Message}, nil
}

func (x *Session) Close() error {
	return nil
}
