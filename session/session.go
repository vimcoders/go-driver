package session

import (
	"bufio"
	"context"
	"errors"
	"go-driver/log"
	"net"
	"runtime/debug"
	"time"
)

type Session struct {
	net.Conn
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
		if err := x.Handle(ctx, request); err != nil {
			return err
		}
	}
}

func (x *Session) Push(ctx context.Context, response Response) (int, error) {
	return x.Write(response)
}

func (x *Session) Close() error {
	return nil
}
