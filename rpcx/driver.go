package rpcx

import (
	"context"
	"go-driver/pb"
	"net"

	"google.golang.org/protobuf/proto"
)

var HandlerDesc = pb.Handler_ServiceDesc

type HandlerClient pb.HandlerClient

type Client interface {
	Close() error
	RemoteAddr() net.Addr
	Register(any) error
	Keeplive(context.Context) error
	Go(context.Context, string, proto.Message) error
	HandlerClient
}
