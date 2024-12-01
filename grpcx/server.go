package grpcx

import (
	"context"
	"net"

	"github.com/vimcoders/go-driver/driver"
	"github.com/vimcoders/go-driver/log"
)

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(ctx context.Context, listener net.Listener, handler driver.Handler) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		handler.Handle(ctx, conn)
	}
}
