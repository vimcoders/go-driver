package handler

import (
	"context"
	"errors"
	"go-driver/app/comet/driver"
	"go-driver/pb"
	"go-driver/rpcx"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	driver.Marshal
	driver.Unmarshal
	*driver.Session
	iClient *rpcx.Client
	Token   string
}

func (x *Session) Handle(ctx context.Context, req driver.Request) error {
	request, err := x.Unmarshal.Unmarshal(req)
	if err != nil {
		return err
	}
	reply, err := x.Unmarshal.Unmarshal(driver.NewRequest(req.Reply()))
	if err != nil {
		return err
	}
	if len(x.Token) > 0 {
		if err := x.iClient.Call(context.Background(), request, reply, &pb.Option{Key: "token", Value: x.Token}); err != nil {
			return err
		}
		if err := x.Push(ctx, reply); err != nil {
			return err
		}
		return nil
	}
	loginRequest, ok := request.(*pb.LoginRequest)
	if !ok {
		return errors.New("!ok")
	}
	if err := x.iClient.Call(context.Background(), request, reply, &pb.Option{Key: "token", Value: loginRequest.Token}); err != nil {
		return err
	}
	loginResponse, ok := reply.(*pb.LoginResponse)
	if !ok {
		return errors.New("!ok")
	}
	if loginResponse.Code != http.StatusOK {
		return errors.New("!ok")
	}
	x.Token = loginRequest.Token
	if err := x.Push(ctx, reply); err != nil {
		return err
	}
	return nil
}

func (x *Session) Push(ctx context.Context, message proto.Message) error {
	response, err := x.Marshal.Marshal(message)
	if err != nil {
		return err
	}
	if _, err := x.Write(response); err != nil {
		return err
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}
