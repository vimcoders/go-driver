package handler

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"runtime/debug"

	"github.com/vimcoders/go-driver/log"
	"github.com/vimcoders/go-driver/quicx"
)

func (x *Handler) ListenAndServe(ctx context.Context) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
	}()
	addr, err := net.ResolveTCPAddr("tcp4", x.TCP.Internet)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	// tcpAddr := listener.Addr().(*net.TCPAddr)
	// listener, err := quicx.Listen("udp", opt.Addr.Port, GenerateTLSConfig(), &quicx.Config{
	// 	MaxIdleTimeout: time.Minute,
	// })
	if err != nil {
		panic(err)
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		go quicx.ListenAndServe(ctx, listener, x)
	}
	log.Infof("running %s", listener.Addr().String())
}
