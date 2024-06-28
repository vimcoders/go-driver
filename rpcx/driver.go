package rpcx

import (
	"go-driver/driver"

	"google.golang.org/protobuf/proto"
)

type ResponsePusher = driver.ResponsePusher

type Handler interface {
	ServeRPCX(ResponsePusher, proto.Message) error
}
