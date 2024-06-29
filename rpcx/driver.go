package rpcx

import (
	"context"
	"go-driver/driver"
	"net"

	"google.golang.org/protobuf/proto"
)

var messages = driver.Messages

type ResponsePusher = driver.ResponsePusher

type Handler interface {
}

type Client interface {
	Call(context.Context, proto.Message, proto.Message) error
	Go(context.Context, proto.Message) error
	Close() error
	RemoteAddr() net.Addr
	Ping(ctx context.Context) (err error)
	Register(h Handler)
	Keeplive(ctx context.Context) error
}
