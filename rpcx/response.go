package rpcx

import (
	"github.com/vimcoders/go-driver/pb"
)

type Response struct {
	RequestId string
	Message   []byte
}

func (x *Response) ToMessage() *pb.Response {
	return &pb.Response{Message: x.Message, RequestId: x.RequestId}
}
