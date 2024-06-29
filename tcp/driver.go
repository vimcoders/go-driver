package tcp

import (
	"context"
	"go-driver/driver"

	"google.golang.org/protobuf/proto"
)

var messages = driver.Messages

type Client = driver.Client

type SHandler interface {
	Handle(ctx context.Context, req, reply proto.Message) error
}

type CHandler interface {
	Handle(ctx context.Context, reply proto.Message) error
}
