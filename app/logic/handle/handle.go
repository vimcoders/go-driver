package handle

import (
	"context"
	"errors"
	"math"
	"net"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"go-driver/app/logic/driver"
	"go-driver/conf"
	"go-driver/log"
	"go-driver/mongox"
	"go-driver/pb"
	"go-driver/rpcx"

	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

var handler = &Handle{}

type Handle struct {
	*mongox.Mongo
	*etcd.Client
	sync.RWMutex
	Users []*driver.User
	Opt   *conf.Conf
	total uint64
	unix  int64
	c     rpcx.Client
}

// MakeHandler creates a Handler instance
func MakeHandler(opt *conf.Conf) *Handle {
	log.Info("etcd endpoints:", opt.Etcd.Endpoints)
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{opt.Etcd.Endpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	log.Info(opt.Mongo.Host, opt.Mongo.DB)
	mongo, err := mongox.Connect(opt.Mongo.Host, opt.Mongo.DB)
	if err != nil {
		panic(err)
	}
	handler.Mongo = mongo
	handler.Client = cli
	handler.Opt = opt
	return handler
}

// Handle receives and executes redis commands
func (x *Handle) Handle(ctx context.Context, conn net.Conn) {
	log.Infof("new conn %s", conn.RemoteAddr().String())
	cli := rpcx.NewClient(conn, math.MaxUint16)
	if err := cli.Register(x); err != nil {
		log.Error(err.Error())
	}
	x.unix = time.Now().Unix()
	//go cli.Keeplive(ctx)
	x.c = cli
	//go cli.Keeplive(context.Background())
}

func (x *Handle) PingRequest(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	x.Lock()
	defer x.Unlock()
	x.total++
	if x.total%500000 == 0 {
		log.Debug("PingRequest", x.total)
	}
	return &pb.PingResponse{}, nil
}

func (x *Handle) Call(ctx context.Context, message proto.Message) (proto.Message, error) {
	// methodName := proto.MessageName(message).Name()
	// method := reflect.ValueOf(x).MethodByName(string(methodName))
	// if ok := method.IsValid(); !ok {
	// 	return nil, errors.New("method.IsValid(); !ok")
	// }
	// args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(message)}
	// result := method.Call(args)
	// if len(result) <= 0 {
	// 	return nil, errors.New("len(result) <= 0")
	// }
	return &pb.PingResponse{}, nil
}

func (x *Handle) Go(ctx context.Context, message proto.Message) error {
	// methodName := proto.MessageName(message).Name()
	// method := reflect.ValueOf(x).MethodByName(string(methodName))
	// if ok := method.IsValid(); !ok {
	// 	return errors.New("method.IsValid(); !ok")
	// }
	// args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(message)}
	// method.Call(args)
	x.c.Go(ctx, &pb.PingResponse{})
	x.Lock()
	defer x.Unlock()
	unix := time.Now().Unix()
	x.total++
	if unix != x.unix {
		log.Debug(x.total, " request/s")
		x.total = 0
		x.unix = unix
	}
	return nil
}

// Close stops handler
func (x *Handle) Close() error {
	for i := 0; i < len(x.Users); i++ {
		// x.Mongo.Insert(x.Users[i])
		// x.Mongo.Update(x.Users[i])
	}
	return nil
}

func (x *Handle) ServeRPCX(w driver.ResponsePusher, in proto.Message) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
		}
		if err != nil {
			log.Error(err.Error())
			debug.PrintStack()
		}
		x.Close()
	}()
	methodName := string(proto.MessageName(in).Name())
	method := reflect.ValueOf(x).MethodByName(methodName)
	values := method.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(in)})
	if len(values) <= 0 {
		return errors.New("len(values) <= 0")
	}
	w.Push(context.Background(), values[0].Interface().(proto.Message))
	return nil
}
