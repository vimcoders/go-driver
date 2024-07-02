package handle

import (
	"context"
	"errors"
	"go-driver/grpcx"
	"go-driver/tcp"
	"reflect"
	"strings"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	tcpclient tcp.Client
	rpcclient grpcx.Client
	Token     string
}

func (x *Session) ServeTCP(ctx context.Context, request proto.Message) error {
	messageName := string(proto.MessageName(request).Name())
	methodName := strings.TrimSuffix(messageName, "Request")
	method := reflect.ValueOf(x.rpcclient).MethodByName(methodName)
	result := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(request)})
	if len(result) <= 0 {
		return errors.New("len(result) <= 0")
	}
	if err := x.tcpclient.Go(ctx, result[0].Interface().(proto.Message)); err != nil {
		return err
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}
