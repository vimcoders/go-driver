package grpcx

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Pusher struct {
	net.Conn
	seq     uint32
	ack     uint32
	timeout time.Duration
	desc    grpc.ServiceDesc
}

func (x *Pusher) Push(_ context.Context, iMessage proto.Message) error {
	messageName := string(proto.MessageName(iMessage).Name())
	methodName := strings.TrimSuffix(messageName, "Request")
	for i := uint16(0); i < uint16(len(x.desc.Methods)); i++ {
		if methodName != x.desc.Methods[i].MethodName {
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
	return fmt.Errorf("%s not registed", methodName)
}
