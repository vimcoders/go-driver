package tcpx

import (
	"bufio"
	"context"
	"errors"
	"go-driver/log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Client interface {
	Register(any) error
	Go(context.Context, proto.Message) error
	net.Conn
	Close() error
}

type Option struct {
	buffsize uint16
	timeout  time.Duration
	grpc.ServiceDesc
}

type ClientX struct {
	net.Conn
	Option
	handler any
}

func NewClient(c net.Conn, opt Option) Client {
	return newClient(c, opt)
}

func newClient(c net.Conn, opt Option) Client {
	opt.buffsize = 8 * 1024
	opt.timeout = time.Second * 60
	x := &ClientX{
		Option: opt,
		Conn:   c,
	}
	return x
}

func (x *ClientX) Go(ctx context.Context, req proto.Message) error {
	return nil
}

func (x *ClientX) Register(a any) error {
	if x.handler != nil {
		return errors.New("x.svr  != nil")
	}
	x.handler = a
	go x.serve(context.Background())
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
	buf := bufio.NewReaderSize(x.Conn, int(x.buffsize))
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.Conn.SetReadDeadline(time.Now().Add(x.timeout)); err != nil {
			return err
		}
		iMessage, err := decode(buf)
		if err != nil {
			return err
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
		if _, err := x.Write(b); err != nil {
			return err
		}
	}
}
