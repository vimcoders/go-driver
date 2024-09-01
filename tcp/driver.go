package tcp

import (
	"context"
	"net"
	"time"
)

type Client interface {
	net.Conn
	Go(context.Context, []byte) error
	Register(Handler) error
}

type Handler interface {
	ServeTCP(context.Context, []byte) error
}

type Option struct {
	Buffsize uint16
	Timeout  time.Duration
}
