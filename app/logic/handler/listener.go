package handler

import (
	"context"
	"encoding/json"
	"go-driver/etcdx"
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
	listener, err := quicx.Listen("udp", x.QUIC.LAN(), GenerateTLSConfig(), &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	go grpcx.ListenAndServe(ctx, listener, x)
	b, err := json.Marshal(&etcdx.Service{WAN: x.QUIC.WAN(), LAN: x.QUIC.LAN()})
	if err != nil {
		panic(err)
	}
	if _, err := x.Client.Put(ctx, x.Etcd.Join("logic"), string(b)); err != nil {
		panic(err)
	}
}
