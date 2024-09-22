package grpcx

import (
	"context"
	"fmt"
	"go-driver/grpcx"
	"go-driver/pb"
	"runtime"
	"time"
)

type Handle struct {
	total uint64
	unix  int64
	pb.ParkourClient
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handle {
	return &Handle{}
}

func (x *Handle) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	unix := time.Now().Unix()
	x.total++
	if unix != x.unix {
		fmt.Println(x.total, " request/s", " NumGoroutine ", runtime.NumGoroutine())
		x.total = 0
		x.unix = unix
	}
	return &pb.PingResponse{Message: req.Message}, nil
}

type Client struct {
	pb.ParkourClient
}

func Dial(network string, addr string) (*Client, error) {
	cli, err := grpcx.Dial(network, addr, grpcx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	if err != nil {
		return nil, err
	}
	cli.Register(context.Background(), MakeHandler())
	return &Client{ParkourClient: pb.NewParkourClient(cli)}, nil
}

func (x *Client) BenchmarkQUIC() {
	for {
		if _, err := x.Ping(context.Background(), &pb.PingRequest{}); err != nil {
			panic(err)
		}
		// if err := x.Go(context.Background(), "Ping", &pb.PingRequest{}); err != nil {
		// 	panic(err)
		// }
	}
}
