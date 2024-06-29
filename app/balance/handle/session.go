package handle

import (
	"context"
	"go-driver/pb"
	"go-driver/rpcx"
	"go-driver/tcp"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	tcpclient tcp.Client
	rpcclient rpcx.Client
	Token     string
}

func (x *Session) ServeTCP(ctx context.Context, request proto.Message) error {
	if len(x.Token) <= 0 {
		return x.Login(ctx, request)
	}
	reply, err := x.rpcclient.Call(context.Background(), request)
	if err != nil {
		return err
	}
	if err := x.tcpclient.Go(ctx, reply); err != nil {
		return err
	}
	return nil
}

func (x *Session) Login(ctx context.Context, request proto.Message) error {
	reply, err := x.rpcclient.Call(context.Background(), request)
	if err != nil {
		return err
	}
	if loginRequest, ok := request.(*pb.LoginRequest); ok {
		x.Token = loginRequest.Token
	}
	if err := x.tcpclient.Go(ctx, reply); err != nil {
		return err
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}
