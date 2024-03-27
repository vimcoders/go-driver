package rpcx

import (
	"context"
	"net"

	"github.com/vimcoders/go-driver/log"

	"github.com/vimcoders/go-driver/driver"
)

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(ctx context.Context, listener net.Listener, handler driver.Handler) {
	// listener, err := net.Listen("tcp", address)
	// if err != nil {
	// 	panic(err)
	// }
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
		log.Debug("new conn", conn.RemoteAddr().String())
		handler.Handle(ctx, conn)
	}
}
