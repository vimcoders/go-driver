package tcp

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"go-driver/log"
	"net"
	"time"
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
	go x.serve(context.Background())
	return nil
}

func (x *XClient) Go(ctx context.Context, request []byte) (err error) {
	if err := x.SetWriteDeadline(time.Now().Add(x.Timeout)); err != nil {
		return err
	}
	if _, err := x.Write(request); err != nil {
		return err
	}
	return nil
}

func (x *XClient) Close() error {
	return x.Conn.Close()
}

func (x *XClient) serve(ctx context.Context) (err error) {
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
		iMessage, err := x.decode(buf)
		if err != nil {
			return err
		}
		if err := x.Handler.ServeTCP(ctx, iMessage); err != nil {
			return err
		}
		if _, err := buf.Discard(len(iMessage)); err != nil {
			return err
		}
	}
}

// 数据流解密
func (x *XClient) decode(b *bufio.Reader) ([]byte, error) {
	headerBytes, err := b.Peek(2)
	if err != nil {
		return nil, err
	}
	header := binary.BigEndian.Uint16(headerBytes)
	if int(header) > b.Size() {
		return nil, fmt.Errorf("header %v too long", header)
	}
	request, err := b.Peek(int(header))
	if err != nil {
		return nil, err
	}
	return request, nil
}
