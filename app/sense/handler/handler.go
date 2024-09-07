package handler

import (
	"context"
	"net"

	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/tcp"

	etcd "go.etcd.io/etcd/client/v3"
)

var handler *Handler

type Handler struct {
	Option
	iClient grpcx.Client
	rpc     grpcx.Client
	pb.UnimplementedParkourServer
	*etcd.Client
}

// MakeHandler creates a Handler instance
func MakeHandler(ctx context.Context) *Handler {
	h := &Handler{}
	if err := h.Parse(); err != nil {
		panic(err)
	}
	if err := h.Connect(ctx); err != nil {
		panic(err)
	}
	handler = h
	return handler
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	newSession := &Session{
		Client:  tcp.NewClient(conn, tcp.Option{}),
		iClient: x.iClient,
	}
	if err := newSession.Register(newSession); err != nil {
		log.Error(err.Error())
	}
}

func (x *Handler) LoginRequest() {
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}
