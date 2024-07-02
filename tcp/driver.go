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
	Ping(context.Context) (err error)
	Register(Handler) error
	Keeplive(context.Context) error
}

type Handler interface {
	ServeTCP(context.Context, proto.Message) error
}
