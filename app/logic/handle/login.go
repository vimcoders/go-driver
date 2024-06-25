package handle

import (
	"net/http"

	"go-driver/driver"
	"go-driver/pb"
	"go-driver/rpcx"
	"go-driver/token"
)

func (x *Handle) LoginRequest(ctx *Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	//ctx.Insert()
	return &pb.LoginResponse{Code: http.StatusOK}, nil
}

func (x *Handle) Authentication(opt rpcx.Option) *Context {
	jwtToken, err := token.ParseToken(opt.Get("token"), []byte(x.Opt.Token.Key))
	if err != nil {
		return &Context{User: &driver.User{}, Mongo: x.Mongo}
	}
	// if user := x.GetUser(jwtToken.Id); user != nil {
	// 	return &Context{User: user}
	// }
	return &Context{User: &driver.User{UserId: jwtToken.Id}, Mongo: x.Mongo}
}
