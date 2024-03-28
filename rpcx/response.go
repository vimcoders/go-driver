package rpcx

import (
	"github.com/vimcoders/go-driver/pb"
)

type ResponseWriter interface {
	Write([]byte) (int, error)
}

type Response struct {
	RequestId string
	Message   []byte
}

func (x *Response) ToMessage() *pb.Response {
	return &pb.Response{Message: x.Message, RequestId: x.RequestId}
}

func (x *Response) Write(b []byte) (int, error) {
	x.Message = b
	return len(b), nil
}
