package handle

import (
	"context"
	"net"
	"runtime"
	"time"

	"go-driver/conf"
	"go-driver/etcdx"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/tcp"

	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

var handler = &Handle{}

type Handle struct {
	rpcclient grpcx.Client
	*etcd.Client
	*conf.Conf
	total uint64
	unix  int64
	pb.UnimplementedHandlerServer
}

// MakeHandler creates a Handler instance
func MakeHandler(opt conf.Conf) *Handle {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{opt.Etcd.Endpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	handler.Conf = &opt
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
		cli, err := grpcx.Dial("udp", response[i].Addr)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if err := cli.Register(x); err != nil {
			log.Error(err.Error())
			continue
		}
		for i := 0; i < 1; i++ {
			go cli.Keeplive(context.Background())
		}
		x.rpcclient = cli
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
	return nil, nil
}

func (x *Handle) Call(ctx context.Context, message proto.Message) (proto.Message, error) {
	return nil, nil
}

func (x *Handle) Go(ctx context.Context, message proto.Message) error {
	// methodName := proto.MessageName(message).Name()
	// method := reflect.ValueOf(x).MethodByName(string(methodName))
	// if ok := method.IsValid(); !ok {
	// 	return errors.New("method.IsValid(); !ok")
	// }
	// args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(message)}
	// method.Call(args)
	unix := time.Now().Unix()
	x.total++
	if unix != x.unix {
		log.Debug(x.total, " request/s", " NumGoroutine ", runtime.NumGoroutine())
		x.total = 0
		x.unix = unix
	}
	return nil
}

// Handle receives and executes redis commands
func (x *Handle) Handle(ctx context.Context, c net.Conn) {
	newSession := &Session{
		tcpclient: tcp.NewClient(c),
		rpcclient: x.rpcclient,
	}
	newSession.tcpclient.Register(newSession)
}

// Close stops handler
func (x *Handle) Close() error {
	return nil
}
