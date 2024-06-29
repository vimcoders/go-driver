package rpcx

import (
	"context"
	"go-driver/driver"

	"google.golang.org/protobuf/proto"
)

var messages = driver.Messages

type ResponsePusher = driver.ResponsePusher

type Client = driver.Client

type Handler interface {
	Call(context.Context, proto.Message) (proto.Message, error)
	Go(context.Context, proto.Message) error
}
