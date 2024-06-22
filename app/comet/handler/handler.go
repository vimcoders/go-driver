package handler

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"go-driver/app/comet/driver"
	"go-driver/conf"
	"go-driver/etcdx"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/quicx"
	"go-driver/rpcx"

	etcd "go.etcd.io/etcd/client/v3"
)

var handler = &Handler{}

type Handler struct {
	driver.Marshal
	driver.Unmarshal
	rpc *rpcx.Client
	*etcd.Client
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
	handler.Conf = &opt
	handler.Client = cli
	handler.Marshal = driver.Messages
	handler.Unmarshal = driver.Messages
	if err := handler.DialLogic(); err != nil {
		panic(err.Error())
	}
	return handler
}

func (x *Handler) DialLogic() error {
	key := x.Etcd.Version + "/service/logic"
	response, err := etcdx.WithQuery[etcdx.Service](x.Client).Query(key)
	if err != nil {
		return err
	}
	for i := 0; i < len(response); i++ {
		log.Info(response[i].Addr)
		conn, err := quicx.Dial(response[i].Addr, &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"quic-echo-example"},
			MaxVersion:         tls.VersionTLS13,
		}, &quicx.Config{
			MaxIdleTimeout: time.Minute,
		})
		//conn, err := net.Dial("tcp", response[i].Addr)
		if err != nil {
			return err
		}
		log.Info(conn.RemoteAddr().String())
		handler.rpc = rpcx.NewClient(conn)
	}
	return nil
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	newSession := &Session{
		rpc:       x.rpc,
		Marshal:   x.Marshal,
		Unmarshal: x.Unmarshal,
		Session: &driver.Session{
			Timeout:  time.Minute * 2,
			Buffsize: 512,
			Conn:     conn,
		},
	}
	newSession.Handler = newSession
	go newSession.Poll(ctx)
}

func (x *Handler) LoginRequest() {
	var replay pb.LoginResponse
	if err := x.rpc.Call(context.Background(), &pb.LoginRequest{}, &replay); err != nil {
		log.Error(err.Error())
	}
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}
