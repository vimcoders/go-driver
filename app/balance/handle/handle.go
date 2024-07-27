package handle

import (
	"context"
	"net"
	"runtime"
	"time"

	"go-driver/app/balance/driver"
	"go-driver/etcdx"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/tcp"

	etcd "go.etcd.io/etcd/client/v3"
)

var handler = &Handle{}
var _ pb.HandlerServer = handler

type Handle struct {
	rpc grpcx.Client
	pb.UnimplementedHandlerServer
	*etcd.Client
	*driver.Option
	total uint64
	unix  int64
}

// MakeHandler creates a Handler instance
func MakeHandler(opt *driver.Option) *Handle {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{opt.Etcd.Endpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	handler.Option = opt
	handler.Client = cli
	if err := handler.DialLogic(); err != nil {
		panic(err.Error())
	}
	return handler
}

func (x *Handle) DialLogic() error {
	key := x.Etcd.Version + "/service/logic"
	response, err := etcdx.WithQuery[etcdx.Service](x.Client).Query(key)
	if err != nil {
		return err
	}
	for i := 0; i < len(response); i++ {
		log.Info(response[i].Addr)
		cli, err := grpcx.Dial("udp", response[i].Addr, pb.Handler_ServiceDesc)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if err := cli.Register(x); err != nil {
			log.Error(err.Error())
			continue
		}
		for i := 0; i < 1; i++ {
			go cli.Keeplive(context.Background(), &pb.PingRequest{})
		}
		cli.Register(x)
		x.rpc = cli
		return nil
	}
	return nil
}

func (x *Handle) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
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
func (x *Handle) Handle(ctx context.Context, c net.Conn) {
	newSession := &Session{
		Client: tcp.NewClient(c),
		rpc:    x.rpc,
	}
	newSession.Client.Register(newSession)
}

// Close stops handler
func (x *Handle) Close() error {
	return nil
}
