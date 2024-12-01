package quicx

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/vimcoders/go-driver/driver"
	"github.com/vimcoders/go-driver/log"
)

// // Config stores tcp server properties
// type Config struct {
// 	Address    string        `yaml:"address"`
// 	MaxConnect uint32        `yaml:"max-connect"`
// 	Timeout    time.Duration `yaml:"timeout"`
// 	Key        string        `yaml:"key"`
// }

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(ctx context.Context, listener net.Listener, handler driver.Handler) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
	}()
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
		log.Debugf("new conn %s", conn.RemoteAddr().String())
		handler.Handle(ctx, conn) // 是否加密
	}
}
