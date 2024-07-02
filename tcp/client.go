package tcp

import (
	"bufio"
	"context"
	"errors"
	"go-driver/log"
	"go-driver/pb"
	"net"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

type XClient struct {
	Handler interface{}
	net.Conn
	sync.RWMutex
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

func (x *XClient) Register(h interface{}) error {
	if x.Handler != nil {
		return errors.New("x.Handler != nil")
	}
	x.Handler = h
	go x.pull(context.Background())
	return nil
}

func (x *XClient) Keeplive(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond)
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
	ticker := time.NewTicker(time.Millisecond * 100)
	buffer := bufio.NewReaderSize(x.Conn, int(x.Buffsize))
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		case <-ticker.C:
			if err := x.Conn.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
				return err
			}
			message, err := decode(buffer)
			if err != nil {
				return err
			}
			x.handle(ctx, message)
		}
	}
}

func (x *XClient) handle(ctx context.Context, message Message) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	req, err := x.new(message.kind())
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(message.message(), req); err != nil {
		return err
	}
	if h, ok := x.Handler.(Handler); ok {
		if err := h.ServeTCP(ctx, req); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (x *XClient) new(kind uint16) (proto.Message, error) {
	if kind >= uint16(len(x.messages)) {
		return nil, errors.New("kind >= uint16(len(x.messages))")
	}
	return x.messages[kind].ProtoReflect().New().Interface(), nil
}
