package handler

import (
	"context"
	"go-driver/driver"
	"go-driver/grpcx"
	"go-driver/tcp"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	tcp.Client
	rpc        grpcx.Client
	Token      string
	MethodDesc driver.MethodDescList
	sync.Pool
}

func (x *Session) ServeTCP(ctx context.Context, stream []byte) error {
	var request driver.Message = stream
	seq := request.Method()
	args := x.MethodDesc[seq].Args
	reply := x.MethodDesc[seq].Replay
	methodName := x.MethodDesc[seq].MethodName
	if err := proto.Unmarshal(request.Payload(), args); err != nil {
		return err
	}
	if err := x.rpc.Invoke(ctx, methodName, args, reply); err != nil {
		return err
	}
	response, err := x.encode(seq, reply)
	if err != nil {
		return err
	}
	if _, err := x.Write(response); err != nil {
		return err
	}
	response.Reset()
	x.Put(&response)
	return nil
}

func (x *Session) Close() error {
	return nil
}

// 数据流加密
func (x *Session) encode(seq uint16, message proto.Message) (driver.Message, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	iMessage := x.Pool.Get().(*driver.Message)
	iMessage.WriteUint16(uint16(4 + len(b)))
	iMessage.WriteUint16(seq)
	if _, err := iMessage.Write(b); err != nil {
		return nil, err
	}
	return *iMessage, nil
}
