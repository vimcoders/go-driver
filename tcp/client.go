package tcp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
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
	for {
		if err := x.Ping(ctx); err != nil {
			log.Error(err.Error())
			return err
		}
	}
}

func (x *XClient) Ping(ctx context.Context) error {
	if err := x.Go(ctx, &pb.PingRequest{}); err != nil {
		return err
	}
	return nil
}

func (x *XClient) Go(ctx context.Context, request proto.Message) (err error) {
	return x.push(context.Background(), request)
}

func (x *XClient) Close() error {
	return x.Conn.Close()
}

func (x *XClient) push(_ context.Context, iMessage proto.Message) error {
	messageName := proto.MessageName(iMessage).Name()
	for i := uint16(0); i < uint16(len(x.messages)); i++ {
		if messageName != proto.MessageName(x.messages[i]).Name() {
			continue
		}
		b, err := encode(i, iMessage)
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
	return fmt.Errorf("%s not registered", messageName)
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
			req := x.messages[iMessage.kind()].ProtoReflect().New().Interface()
			if err := proto.Unmarshal(iMessage.payload(), req); err != nil {
				return fmt.Errorf("%s %v", err.Error(), iMessage)
			}
			if err := x.Handler.ServeTCP(ctx, req); err != nil {
				return err
			}
		}
	}
}
