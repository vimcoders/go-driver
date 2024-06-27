package handle

import (
	"context"
	"go-driver/app/balance/driver"
	"go-driver/pb"
	"go-driver/rpcx"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	h *driver.Handle
	*rpcx.Client
	Token string
}

func (x *Session) Handle(ctx context.Context, request, reply proto.Message) error {
	if len(x.Token) <= 0 {
		return x.Login(ctx, request, reply)
	}
	if err := x.Call(context.Background(), request, reply); err != nil {
		return err
	}
	if err := x.Push(ctx, reply); err != nil {
		return err
	}
	return nil
}

func (x *Session) Login(ctx context.Context, request, reply proto.Message) error {
	if err := x.Call(context.Background(), request, reply); err != nil {
		return err
	}
	if loginRequest, ok := request.(*pb.LoginRequest); ok {
		x.Token = loginRequest.Token
	}
	if err := x.Push(ctx, reply); err != nil {
		return err
	}
	return nil
}

func (x *Session) Push(ctx context.Context, message proto.Message) error {
	if err := x.h.Push(ctx, message); err != nil {
		return err
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}
