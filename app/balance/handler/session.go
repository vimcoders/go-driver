package handler

import (
	"context"
	"fmt"
	"go-driver/driver"
	"go-driver/grpcx"
	"go-driver/tcp"
	"strings"

	"google.golang.org/protobuf/proto"
)

var Messages = driver.Messages

type Session struct {
	tcp.Client
	rpc   grpcx.Client
	Token string
}

func (x *Session) ServeTCP(ctx context.Context, request proto.Message) error {
	requestName := string(proto.MessageName(request).Name())
	methodName := strings.TrimSuffix(requestName, "Request")
	responseName := methodName + "Response"
	for i := 0; i < len(Messages); i++ {
		if string(proto.MessageName(Messages[i]).Name()) != responseName {
			continue
		}
		reply := Messages[i].ProtoReflect().New().Interface()
		if err := x.rpc.Invoke(ctx, methodName, request, reply); err != nil {
			return err
		}
		if err := x.Go(ctx, reply); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("%s not in Messages", responseName)
}

func (x *Session) Close() error {
	return nil
}
