package handle

import (
	"context"
	"errors"
	"go-driver/grpcx"
	"go-driver/tcp"
	"reflect"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	tcpclient tcp.Client
	rpcclient grpcx.Client
	Token     string
}

func (x *Session) ServeTCP(ctx context.Context, request proto.Message) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	messageName := string(proto.MessageName(request).Name())
	methodName := strings.TrimSuffix(messageName, "Request")
	method := reflect.ValueOf(x.rpcclient).MethodByName(methodName)
	result := method.Call([]reflect.Value{reflect.ValueOf(timeoutCtx), reflect.ValueOf(request)})
	if len(result) <= 0 {
		return errors.New("len(result) <= 0")
	}
	x.tcpclient.Go(ctx, result[0].Interface().(proto.Message))
	return nil
}

func (x *Session) Close() error {
	return nil
}
