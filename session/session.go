package session

import (
	"bufio"
	"context"
	"errors"
	"go-driver/driver"
	"go-driver/log"
	"net"
	"runtime/debug"
	"time"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	net.Conn
	Id    string
	Token string
	driver.Marshal
	driver.Unmarshal
	Buffsize int
	Timeout  time.Duration
	Handler
}

func (x *Session) Poll(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
		}
		if err != nil {
			log.Error(err.Error())
			debug.PrintStack()
		}
		x.Close()
	}()
	buffer := bufio.NewReaderSize(x.Conn, x.Buffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		request, err := ReadRequest(buffer)
		if err != nil {
			return err
		}
		if err := x.Handle(x, request); err != nil {
			return err
		}
		if _, err := buffer.Discard(len(request)); err != nil {
			return err
		}
	}
}

func (x *Session) Push(ctx context.Context, message proto.Message) (int, error) {
	b, err := x.Marshal.Marshal(message)
	if err != nil {
		return 0, err
	}
	return x.Write(b)
}

func (x *Session) Handle(w driver.ResponsePusher, req Request) error {
	request, err := x.Unmarshal.Unmarshal(req)
	if err != nil {
		return err
	}
	reply, err := x.Unmarshal.Unmarshal(NewRequest(req.Kind() + 1))
	if err != nil {
		return err
	}
	if x.Handler != nil {
		return x.Handler.Handle(w, request, reply)
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}
