package kcpx

import (
	"context"
	"errors"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type Client interface {
	Go(context.Context, proto.Message) error
	Close() error
	RemoteAddr() net.Addr
	Register(Handler) error
	Keeplive(context.Context, proto.Message) error
}

type Handler interface {
	ServeKCP(context.Context, proto.Message) error
}

type Option struct {
	Messages []proto.Message
	Buffsize uint16
	Timeout  time.Duration
}

func (x *Option) encode(m proto.Message) (Message, error) {
	messageName := string(proto.MessageName(m).Name())
	for i := 0; i < len(x.Messages); i++ {
		if string(proto.MessageName(x.Messages[i]).Name()) != messageName {
			continue
		}
		return encode(uint16(i), m)
	}
	return nil, nil
}

func (x *Option) decode(m Message) (proto.Message, error) {
	reqestId := m.req()
	if int(reqestId) >= len(x.Messages) {
		return nil, errors.New("messageId >= len(x.messages)")
	}
	newMessage := x.Messages[reqestId].ProtoReflect().New().Interface()
	if err := proto.Unmarshal(m.payload(), newMessage); err != nil {
		return nil, err
	}
	return newMessage, nil
}
