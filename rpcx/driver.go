package rpcx

import (
	"context"
	"go-driver/driver"

	"google.golang.org/protobuf/proto"
)

var messages = driver.Messages

type ResponsePusher = driver.ResponsePusher

type Handler interface {
}

type Client interface {
	Call(context.Context, proto.Message, proto.Message) error
	Go(context.Context, proto.Message) error
	Keeplive(context.Context) error
	Close() error
	RemoteAddr() string
	Ping(ctx context.Context) (err error)
}
