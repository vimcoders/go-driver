package driver

import (
	"context"
	"net"

	"google.golang.org/protobuf/proto"
)

type Client interface {
	Call(context.Context, proto.Message) (proto.Message, error)
	Go(context.Context, proto.Message) error
	Close() error
	RemoteAddr() net.Addr
	Ping(ctx context.Context) (err error)
	Register(interface{}) error
	Keeplive(ctx context.Context) error
}
