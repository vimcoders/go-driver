package tcp

import (
	"bufio"
	"context"
	"errors"
	"go-driver/log"
	"go-driver/pb"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type XClient struct {
	Handler Handler
	net.Conn
	Buffsize uint16
	Timeout  time.Duration
	messages []proto.Message
}

func NewClient(c net.Conn) Client {
	x := &XClient{
		Conn:     c,
		Buffsize: 16 * 1024,
		Timeout:  time.Second * 240,
		messages: messages,
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

func (x *XClient) Keeplive(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond * 1)
	for range ticker.C {
		if err := x.Ping(ctx); err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (x *XClient) Ping(ctx context.Context) error {
	if err := x.Go(ctx, &pb.PingRequest{}); err != nil {
		return err
	}
	return nil
}

func (x *XClient) Go(ctx context.Context, request proto.Message) (err error) {
	pusher := &Pusher{
		Conn:     x.Conn,
		timeout:  x.Timeout,
		messages: x.messages,
	}
	return pusher.Push(context.Background(), request)
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
			if err := x.Conn.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
				return err
			}
			iMessage, err := decode(buf)
			if err != nil {
				return err
			}
			req, err := x.new(iMessage.kind())
			if err != nil {
				return err
			}
			if err := proto.Unmarshal(iMessage.message(), req); err != nil {
				return err
			}
			if err := x.Handler.ServeTCP(ctx, req); err != nil {
				return err
			}
			return nil
		}
	}
}

func (x *XClient) new(kind uint16) (proto.Message, error) {
	if kind >= uint16(len(x.messages)) {
		return nil, errors.New("kind >= uint16(len(x.messages))")
	}
	return x.messages[kind].ProtoReflect().New().Interface(), nil
}
