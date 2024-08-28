package handler

import (
	"context"
	"fmt"
	"go-driver/log"
	"go-driver/quicx"
	"net"
	"runtime"
	"runtime/debug"
)

func (x *Handler) ListenAndServe(ctx context.Context) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
	}()
	addr, err := net.ResolveTCPAddr("tcp4", x.TCP.WAN())
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
