package handler

import (
	"context"
	"go-driver/grpcx"
	"go-driver/quicx"
	"time"
)

func (x *Handler) ListenAndServe(ctx context.Context) {
	// addr, err := net.ResolveTCPAddr("tcp4", opt.Addr.Port)
	// if err != nil {
	// 	panic(err)
	// }
	// listener, err := net.ListenTCP("tcp", addr)
	listener, err := quicx.Listen("udp", x.QUIC.Local, GenerateTLSConfig(), &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	go grpcx.ListenAndServe(ctx, listener, x)
}
