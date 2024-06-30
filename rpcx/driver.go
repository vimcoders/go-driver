package rpcx

import (
	"context"
	"go-driver/driver"
	"go-driver/pb"
	"net"

	"google.golang.org/protobuf/proto"
)

var HandlerDesc = pb.Handler_ServiceDesc

type ResponsePusher = driver.ResponsePusher

type HandlerClient pb.HandlerClient

type Client interface {
	Close() error
	RemoteAddr() net.Addr
	Register(any) error
	Keeplive(ctx context.Context) error
	Go(context.Context, string, proto.Message) error
	HandlerClient
}
