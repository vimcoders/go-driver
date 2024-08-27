package handler

import (
	"context"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/pb"
	"net"

	etcd "go.etcd.io/etcd/client/v3"
)

type Handler struct {
	Option
	*etcd.Client
	pb.ChatServer
}

func MakeHandler(ctx context.Context) *Handler {
	h := &Handler{}
	if err := h.Parse(); err != nil {
		panic(err)
	}
	if err := h.Connect(ctx); err != nil {
		panic(err)
	}
	return h
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	log.Infof("new conn %s", conn.RemoteAddr().String())
	cli := grpcx.NewClient(conn, grpcx.Option{ServiceDesc: pb.Chat_ServiceDesc})
	if err := cli.Register(x); err != nil {
		log.Error(err.Error())
	}
	for i := 0; i < 1; i++ {
		go cli.Keeplive(context.Background(), &pb.PingRequest{})
	}
}
