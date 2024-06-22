package handler

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"go-driver/conf"
	"go-driver/driver"
	"go-driver/etcdx"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/quicx"
	"go-driver/rpcx"
	"go-driver/session"

	etcd "go.etcd.io/etcd/client/v3"
)

var handler = &Handler{}

type Handler struct {
	iClient *rpcx.Client
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
	response, err := etcdx.WithQuery[etcdx.Service](cli).Query(key)
	if err != nil {
		panic(err.Error())
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
			panic(err)
		}
		log.Info(conn.RemoteAddr().String())
		handler.iClient = rpcx.NewClient(conn)
	}
	handler.Marshal = session.Messages
	handler.Unmarshal = session.Messages
	return handler
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	newSession := &Session{
		iClient: x.iClient,
		Session: &session.Session{
			Timeout:   time.Minute * 2,
			Buffsize:  512,
			Conn:      conn,
			Marshal:   x.Marshal,
			Unmarshal: x.Unmarshal,
		},
	}
	newSession.Handler = newSession
	go newSession.Poll(ctx)
}

func (x *Handler) LoginRequest() {
	var replay pb.LoginResponse
	if err := x.iClient.Call(context.Background(), &pb.LoginRequest{}, &replay); err != nil {
		log.Error(err.Error())
	}
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}
