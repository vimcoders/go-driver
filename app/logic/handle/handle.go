package handle

import (
	"context"
	"errors"
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
	handle := &rpcx.Handle{
		Conn:     conn,
		Buffsize: 16 * 1024,
		Timeout:  time.Second * 120,
	}
	handle.Register(x)
}

func (x *Handle) PingRequest(ctx *Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	x.Lock()
	defer x.Unlock()
	log.Debug(req)
	// x.total++
	// log.Debugf("%v", x.total)
	return &pb.PingResponse{}, nil
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
