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
	Handler  // handler to invoke, http.DefaultServeMux if nil
	Buffsize uint16
	Timeout  time.Duration
	messages []proto.Message
}

func (x *Handle) Register(handler Handler, p ...proto.Message) error {
	if len(x.messages) > 0 {
		return errors.New("len(x.ProtoBuf) > 0")
	}
	x.Handler = handler
	t, v := reflect.TypeOf(x.Handler), reflect.ValueOf(x.Handler)
	for i := 0; i < t.NumMethod(); i++ {
		method := v.Method(i)
		kind := method.Type()
		if kind.NumIn() < 2 {
			continue
		}
		e := t.In(1).Elem()
		in, ok := reflect.New(e).Interface().(proto.Message)
		if !ok {
			continue
		}
		x.messages = append(x.messages, in)
	}
	x.messages = append(x.messages, p...)
	var message pb.Message
	for i := 0; i < len(x.messages); i++ {
		message.Messages = append(message.Messages, string(proto.MessageName(x.messages[i]).Name()))
	}
	if err := x.Push(context.Background(), &message); err != nil {
		return err
	}
	go x.Pull(context.Background())
	return nil
}

func (x *Handle) serveRPCX(message Message) (err error) {
	request, err := x.Unmarshal(message)
	if err != nil {
		return err
	}
	pusher := &Pusher{
		Conn:     x.Conn,
		Timeout:  x.Timeout,
		messages: x.messages,
		Message:  message,
	}
	x.Handler.ServeRPCX(pusher, request)
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
		go x.serveRPCX(message)
	}
}

func (x *Handle) Unmarshal(request Message) (proto.Message, error) {
	if len(request) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	kind := request.Kind()
	if kind >= uint16(len(x.messages)) {
		return nil, fmt.Errorf("message id %v not registered", kind)
	}
	message := x.messages[kind].ProtoReflect().New().Interface()
	if err := proto.Unmarshal(request.Message(), message); err != nil {
		return nil, err
	}
	return message, nil
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
	Timeout  time.Duration
	messages []proto.Message
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
	for i := uint16(0); i < uint16(len(x.messages)); i++ {
		if proto.MessageName(message) != proto.MessageName(x.messages[i]) {
			continue
		}
		return encode(x.SeqNumber(), i, message)
	}
	return nil, fmt.Errorf("message %s not registered", proto.MessageName(message))
}
