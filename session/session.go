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
	w        net.Conn
	Buffsize int
	Timeout  time.Duration
	Handler
}

func NewSession(w net.Conn) *Session {
	return &Session{w: w, Buffsize: 512, Timeout: time.Minute}
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
	buffer := bufio.NewReaderSize(x.w, x.Buffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.w.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		request, err := parseRequest(buffer)
		if err != nil {
			return err
		}
		if err := x.Handle(ctx, request); err != nil {
			return err
		}
	}
}

func (x *Session) Push(ctx context.Context, response Response) (int, error) {
	return x.w.Write(response)
}

func (x *Session) Close() error {
	return nil
}
