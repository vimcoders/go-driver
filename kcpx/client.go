package kcpx

import (
	"bufio"
	"context"
	"errors"
	"go-driver/log"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type XClient struct {
	Option
	Handler Handler
	net.Conn
}

func NewClient(c net.Conn, opt Option) Client {
	x := &XClient{
		Conn:   c,
		Option: opt,
	}
	if x.Buffsize <= 0 {
		x.Option.Buffsize = 1024
	}
	if x.Timeout <= 0 {
		x.Option.Timeout = time.Second
	}
	return x
}

func (x *XClient) Register(h Handler) error {
	if x.Handler != nil {
		return errors.New("x.Handler != nil")
	}
	x.Handler = h
	go x.pull(context.Background())
	return nil
}

func (x *XClient) Keeplive(ctx context.Context, ping proto.Message) error {
	for {
		if err := x.Go(ctx, ping); err != nil {
			log.Error(err.Error())
			return err
		}
	}
}

func (x *XClient) Go(ctx context.Context, request proto.Message) (err error) {
	b, err := x.encode(request)
	if err != nil {
		return err
	}
	if err := x.SetWriteDeadline(time.Now().Add(x.Timeout)); err != nil {
		return err
	}
	if _, err := b.WriteTo(x.Conn); err != nil {
		return err
	}
	return nil
}

func (x *XClient) Close() error {
	return x.Conn.Close()
}

func (x *XClient) pull(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
		if err := x.Close(); err != nil {
			log.Error(err.Error())
		}
	}()
	buf := bufio.NewReaderSize(x.Conn, int(x.Buffsize))
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.Conn.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		iMessage, err := decode(buf)
		if err != nil {
			return err
		}
		req, err := x.decode(iMessage)
		if err != nil {
			return err
		}
		if err := x.Handler.ServeKCP(ctx, req); err != nil {
			return err
		}
	}
}
