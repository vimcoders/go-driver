package handler

import (
	"context"
	"net/http"

	"go-driver/app/logic/driver"
	"go-driver/log"
	"go-driver/pb"

	"google.golang.org/protobuf/proto"
)

func (x *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Debug(req)
	return &pb.LoginResponse{Code: http.StatusOK}, nil
}

func (x *Handler) Authentication(request proto.Message) *Context {
	// jwtToken, err := token.ParseToken(opt.Get("token"), []byte(x.Opt.Token.Key))
	// if err != nil {
	// 	return &Context{User: &driver.User{}, Mongo: x.Mongo}
	// }
	// if user := x.GetUser(jwtToken.Id); user != nil {
	// 	return &Context{User: user}
	// }
	return &Context{User: &driver.User{}, Mongo: x.Mongo}
}
