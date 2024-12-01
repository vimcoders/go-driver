package handler

import (
	"context"
	"net"
	"runtime"
	"time"

	"github.com/vimcoders/go-driver/pb"

	"github.com/vimcoders/go-driver/udpx"

	"github.com/vimcoders/go-driver/log"
	//etcd "go.etcd.io/etcd/client/v3"
)

var handler *Handler
var _ pb.ParkourServer = handler

type Handler struct {
	//rpc grpcx.Client
	pb.UnimplementedParkourServer
	//*etcd.Client
	Option
	total uint64
	unix  int64
}

// MakeHandler creates a Handler instance
func MakeHandler(ctx context.Context) *Handler {
	h := &Handler{}
	if err := h.Parse(); err != nil {
		panic(err)
	}
	if err := h.Connect(ctx); err != nil {
		panic(err)
	}
	handler = h
	return handler
}

func (x *Handler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	unix := time.Now().Unix()
	x.total++
	if unix != x.unix {
		log.Debug(x.total, " request/s", " NumGoroutine ", runtime.NumGoroutine())
		x.total = 0
		x.unix = unix
	}
	return &pb.PingResponse{Message: req.Message}, nil
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, c net.Conn) {
	newSession := &Session{}
	cli := udpx.NewClient(c.(*net.UDPConn), udpx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	//cli := tcpx.NewClient(c, tcpx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	if err := cli.Register(ctx, newSession); err != nil {
		log.Error(err.Error())
	}
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}
