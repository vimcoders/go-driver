package handle

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"go-driver/app/comet/driver"
	"go-driver/conf"
	"go-driver/etcdx"
	"go-driver/log"
	"go-driver/quicx"
	"go-driver/rpcx"

	etcd "go.etcd.io/etcd/client/v3"
)

var handler = &Handle{}

type Handle struct {
	driver.Marshal
	driver.Unmarshal
	c *rpcx.Client
	*etcd.Client
	*conf.Conf
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
	handler.Marshal = driver.Messages
	handler.Unmarshal = driver.Messages
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
		handler.c = rpcx.NewClient(conn)
	}
	return nil
}

// Handle receives and executes redis commands
func (x *Handle) Handle(ctx context.Context, c net.Conn) {
	newSession := &Session{
		Client:    x.c,
		Marshal:   x.Marshal,
		Unmarshal: x.Unmarshal,
		h:         driver.NewHandle(c),
	}
	newSession.h.Handler = newSession
	go newSession.h.Pull(ctx)
}

// Close stops handler
func (x *Handle) Close() error {
	return nil
}
