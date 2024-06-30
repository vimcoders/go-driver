package handle

import (
	"context"
	"errors"
	"math"
	"net"
	"reflect"
	"runtime"
	"sync"
	"time"

	"go-driver/app/logic/driver"
	"go-driver/conf"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/mongox"
	"go-driver/pb"

	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

var handler = &Handle{}

type Handle struct {
	*mongox.Mongo
	*etcd.Client
	Users []*driver.User
	Opt   *conf.Conf
	total uint64
	unix  int64
	c     grpcx.Client
	sync.RWMutex
	pb.UnimplementedHandlerServer
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
	cli := grpcx.NewClient(conn, math.MaxUint32/2)
	if err := cli.Register(x); err != nil {
		log.Error(err.Error())
	}
	x.unix = time.Now().Unix()
	x.c = cli
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

func (x *Handle) Call(ctx context.Context, message proto.Message) (proto.Message, error) {
	methodName := proto.MessageName(message).Name()
	method := reflect.ValueOf(x).MethodByName(string(methodName))
	if ok := method.IsValid(); !ok {
		return nil, errors.New("method.IsValid(); !ok")
	}
	args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(message)}
	result := method.Call(args)
	if len(result) <= 0 {
		return nil, errors.New("len(result) <= 0")
	}
	return result[0].Interface().(proto.Message), nil
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

// Close stops handler
func (x *Handle) Close() error {
	for i := 0; i < len(x.Users); i++ {
		// x.Mongo.Insert(x.Users[i])
		// x.Mongo.Update(x.Users[i])
	}
	return nil
}
