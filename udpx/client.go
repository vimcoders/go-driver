package udpx

import (
	"context"
	"errors"
	"fmt"
	"go-driver/log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Client interface {
	Register(context.Context, any) error
	Go(context.Context, proto.Message) error
	net.Conn
	Close() error
}

type Handler interface {
	Handle(context.Context, Message) error
}

type Option struct {
	buffsize uint16
	timeout  time.Duration
	grpc.ServiceDesc
}

type ClientX struct {
	*net.UDPConn
	Option
	handler any
}

func NewClient(c *net.UDPConn, opt Option) Client {
	return newClient(c, opt)
}

func newClient(c *net.UDPConn, opt Option) Client {
	opt.buffsize = 8 * 1024
	opt.timeout = time.Second * 60
	x := &ClientX{
		Option:  opt,
		UDPConn: c,
	}
	return x
}

func (x *ClientX) Go(ctx context.Context, req proto.Message) error {
	metodName := proto.MessageName(req)
	// for methodId := 0; methodId < len(x.Methods); methodId++ {
	// 	if ok := strings.Contains(string(metodName), x.Methods[methodId].MethodName); !ok {
	// 		continue
	// 	}
	// 	if err := x.push(uint16(methodId), req); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }
	return fmt.Errorf("%s not registed", metodName)
}

func (x *ClientX) Register(ctx context.Context, a any) error {
	if x.handler != nil {
		return errors.New("x.svr  != nil")
	}
	x.handler = a
	go x.serve(ctx)
	return nil
}

func (x *ClientX) serve(ctx context.Context) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
		if err != nil {
			log.Error(err.Error())
		}
		if err := x.Close(); err != nil {
			log.Error(err.Error())
		}
	}()
	iMessage := make(Message, x.buffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		_, addr, err := x.UDPConn.ReadFromUDP(iMessage)
		if err != nil {
			return err
		}
		if x.handler == nil {
			continue
		}
		method, payload := iMessage.method(), iMessage.payload()
		dec := func(in any) error {
			if err := proto.Unmarshal(payload, in.(proto.Message)); err != nil {
				return err
			}
			return nil
		}
		reply, err := x.Methods[method].Handler(x.handler, ctx, dec, nil)
		if err != nil {
			return err
		}
		b, err := encode(method, reply.(proto.Message))
		if err != nil {
			return err
		}
		if _, err := x.WriteToUDP(b, addr); err != nil {
			log.Error(err)
		}
		b.reset()
	}
}
