package handle

import (
	"context"
	"go-driver/app/balance/driver"
	"go-driver/pb"
	"go-driver/rpcx"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	driver.Marshal
	driver.Unmarshal
	h *driver.Handle
	*rpcx.Client
	Token string
}

func (x *Session) Handle(ctx context.Context, req driver.Request) error {
	if len(x.Token) <= 0 {
		return x.Login(ctx, req)
	}
	request, err := x.Unmarshal.Unmarshal(req)
	if err != nil {
		return err
	}
	reply, err := x.Unmarshal.Unmarshal(driver.NewRequest(req.Reply()))
	if err != nil {
		return err
	}
	if err := x.Call(context.Background(), request, reply, &pb.Option{Key: "token", Value: x.Token}); err != nil {
		return err
	}
	if err := x.Push(ctx, reply); err != nil {
		return err
	}
	return nil
}

func (x *Session) Login(ctx context.Context, req driver.Request) error {
	request, reply := pb.LoginRequest{}, pb.LoginResponse{}
	if err := proto.Unmarshal(req.Message(), &request); err != nil {
		return err
	}
	if err := x.Call(context.Background(), &request, &reply); err != nil {
		return err
	}
	x.Token = request.Token
	if err := x.Push(ctx, &reply); err != nil {
		return err
	}
	return nil
}

func (x *Session) Push(ctx context.Context, message proto.Message) error {
	response, err := x.Marshal.Marshal(message)
	if err != nil {
		return err
	}
	if _, err := x.h.Push(ctx, response); err != nil {
		return err
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}
