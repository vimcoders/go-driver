package message

import (
	"github.com/vimcoders/go-driver/pb"
	"google.golang.org/protobuf/proto"
)

var GateMessages = []proto.Message{
	&pb.LoginRequest{},
	&pb.LoginResponse{},
}
