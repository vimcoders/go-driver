package grpcx

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// 推送一个proto对象到对端
type Pusher struct {
	net.Conn
	seq     uint32
	ack     uint32
	timeout time.Duration
	method  string
	desc    grpc.ServiceDesc
}

func (x *Pusher) Push(_ context.Context, iMessage proto.Message) error {
	for i := uint16(0); i < uint16(len(x.desc.Methods)); i++ {
		if x.method != x.desc.Methods[i].MethodName {
			continue
		}
		b, err := encode(x.seq, x.ack, i, iMessage)
		if err != nil {
			return err
		}
		if err := x.SetWriteDeadline(time.Now().Add(x.timeout)); err != nil {
			return err
		}
		if _, err := x.Conn.Write(b); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("%s not registed", x.method)
}
