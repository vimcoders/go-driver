package tcp

import (
	"context"
	"go-driver/driver"

	"google.golang.org/protobuf/proto"
)

var messages = driver.Messages

type Client = driver.Client

type Handler interface {
	ServeTCP(ctx context.Context, req proto.Message) error
}
