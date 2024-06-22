package handler

import (
	"context"
	"errors"
	"go-driver/app/comet/driver"
	"go-driver/pb"
	"go-driver/rpcx"
	"go-driver/session"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	iClient *rpcx.Client
	*session.Session
}

func (x *Session) Handle(w driver.ResponsePusher, request, reply proto.Message) error {
	if len(x.Token) > 0 {
		if err := x.iClient.Call(context.Background(), request, reply, &pb.Option{Key: "token", Value: x.Token}); err != nil {
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
	return nil
}

func (x *Session) Close() error {
	return nil
}
