package tcp

import (
	"context"
	"go-driver/driver"
	"net"

	"google.golang.org/protobuf/proto"
)

var messages = driver.Messages

type Client interface {
	Go(context.Context, proto.Message) error
	Close() error
	RemoteAddr() net.Addr
	Ping(ctx context.Context) (err error)
	Register(interface{}) error
	Keeplive(ctx context.Context) error
}

type Handler interface {
	ServeTCP(ctx context.Context, req proto.Message) error
}
