package driver

import (
	"go-driver/pb"

	"google.golang.org/protobuf/proto"
)

// 定义所有的协议号
var Messages = []proto.Message{
	&pb.PingRequest{},
	&pb.PingResponse{},
	&pb.LoginRequest{},
	&pb.LoginResponse{},
}
