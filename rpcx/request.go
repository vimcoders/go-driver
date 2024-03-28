package rpcx

import (
	"github.com/vimcoders/go-driver/pb"
)

type Request struct {
	RequestId string
	Message   []byte
}

func (x *Request) ToMessage() *pb.Request {
	return &pb.Request{Message: x.Message, RequestId: x.RequestId}
}
