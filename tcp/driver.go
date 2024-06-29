package tcp

import (
	"context"
	"go-driver/driver"

	"google.golang.org/protobuf/proto"
)

var messages = driver.Messages

type Client = driver.Client

type SHandler interface {
	ServeTCP(ctx context.Context, req, reply proto.Message) error
}

type CHandler interface {
	ServeTCP(ctx context.Context, reply proto.Message) error
}
