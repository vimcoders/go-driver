package rpcx

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go-driver/log"
	"go-driver/pb"
	"net"
	"reflect"
	"runtime/debug"
	"time"

	"google.golang.org/protobuf/proto"
)

type Handle struct {
	net.Conn
	Handler  interface{} // handler to invoke, http.DefaultServeMux if nil
	Buffsize uint16
	Timeout  time.Duration
	ProtoBuf
}

func (x *Handle) Register(handler interface{}, p ...proto.Message) error {
	if len(x.ProtoBuf) > 0 {
		return errors.New("len(x.ProtoBuf) > 0")
	}
	x.Handler = handler
	x.ProtoBuf = p
	go x.Pull(context.Background())
	return nil
}

func (x *Handle) ServeRPCX(w ResponsePusher, in proto.Message) (err error) {
	methodName := string(proto.MessageName(in).Name())
	method := reflect.ValueOf(x.Handler).MethodByName(methodName)
	values := method.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(in)})
	if len(values) <= 0 {
		return errors.New("len(values) <= 0")
	}
	w.Push(context.Background(), values[0].Interface().(proto.Message))
	return nil
}

func (x *Handle) Pull(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
			debug.PrintStack()
		}
		if err := x.Close(); err != nil {
			log.Error(err.Error())
		}
	}()
	buffer := bufio.NewReaderSize(x.Conn, int(x.Buffsize))
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		message, err := decode(buffer)
		if err != nil {
			return err
		}
		request, err := x.ProtoBuf.Unmarshal(message)
		if err != nil {
			return err
		}
		pusher := &Pusher{
			Conn:    x.Conn,
			Timeout: x.Timeout,
		}
		x.ServeRPCX(pusher, request)
	}
}

func (x *Handle) Push(ctx context.Context, iMessage *pb.Message) error {
	pusher := Pusher{
		Conn:    x.Conn,
		Timeout: x.Timeout,
	}
	return pusher.Push(ctx, iMessage)
}

type Pusher struct {
	Message
	net.Conn
	Timeout time.Duration
	ProtoBuf
}

func (x *Pusher) Push(ctx context.Context, iMessage proto.Message) error {
	response, err := x.Marshal(iMessage)
	if err != nil {
		return err
	}
	if err := x.SetWriteDeadline(time.Now().Add(x.Timeout)); err != nil {
		return err
	}
	if _, err := x.Conn.Write(response); err != nil {
		return err
	}
	return nil
}

func (x *Pusher) Marshal(message proto.Message) ([]byte, error) {
	for i := uint16(0); i < uint16(len(x.ProtoBuf)); i++ {
		if proto.MessageName(message) != proto.MessageName(x.ProtoBuf[i]) {
			continue
		}
		return encode(x.PackageNumber(), i, message)
	}
	return nil, fmt.Errorf("message %s not registered", proto.MessageName(message))
}
