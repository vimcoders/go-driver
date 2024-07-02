package handle

import (
	"context"
	"go-driver/grpcx"
	"go-driver/pb"
	"go-driver/tcp"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

var HandlerDesc = pb.Handler_ServiceDesc

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
	dec := func(in any) error {
		b, err := proto.Marshal(request)
		if err != nil {
			return err
		}
		if err := proto.Unmarshal(b, in.(proto.Message)); err != nil {
			return err
		}
		return nil
	}
	for i := 0; i < len(HandlerDesc.Methods); i++ {
		if HandlerDesc.Methods[i].MethodName != methodName {
			continue
		}
		response, err := HandlerDesc.Methods[i].Handler(x.rpcclient, timeoutCtx, dec, nil)
		if err != nil {
			return err
		}
		if err := x.tcpclient.Go(ctx, response.(proto.Message)); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}
