package handler

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"go-driver/conf"
	"go-driver/driver"
	"go-driver/etcdx"
	"go-driver/grpcx"
	"go-driver/log"

	etcd "go.etcd.io/etcd/client/v3"
)

var handler = &Handler{}

type Handler struct {
	iClient *grpcx.Client
	driver.Marshal
	driver.Unmarshal
	*conf.Conf
}

// MakeHandler creates a Handler instance
func MakeHandler(opt conf.Conf) *Handler {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{opt.Etcd.Endpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	key := opt.Etcd.Version + "/service/logic"
	response, err := cli.Get(context.Background(), key, etcd.WithPrefix())
	if err != nil {
		panic(err.Error())
	}
	var service etcdx.Service
	for i := 0; i < len(response.Kvs); i++ {
		log.Info(string(response.Kvs[i].Value))
		if err := json.Unmarshal(response.Kvs[i].Value, &service); err != nil {
			panic(err.Error())
		}
	}
	log.Info(service.Addr)
	// conn, err := quicx.Dial(service.Addr, &tls.Config{
	// 	InsecureSkipVerify: true,
	// 	NextProtos:         []string{"quic-echo-example"},
	// 	MaxVersion:         tls.VersionTLS13,
	// }, &quicx.Config{
	// 	MaxIdleTimeout: time.Minute,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	//message := &driver.Protobuf{Messages: driver.Messages}
	//handler.Marshal = message
	//handler.Unmarshal = message
	//handler.iClient = rpcx.NewClient(conn)
	return handler
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	newSession := &Session{
		iClient:   x.iClient,
		Marshal:   x.Marshal,
		Unmarshal: x.Unmarshal,
		Timeout:   time.Minute * 2,
		Header:    4,
		Buffsize:  1024 * 4,
		Conn:      conn,
	}
	go newSession.Poll(ctx)
}

func (x *Handler) LoginRequest() {
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}
