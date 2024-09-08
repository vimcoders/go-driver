package handler

import (
	"context"
	"net"

	"go-driver/grpcx"
	"go-driver/pb"

	etcd "go.etcd.io/etcd/client/v3"
)

var handler *Handler

type Handler struct {
	Option
	rpc grpcx.Client
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
}

func (x *Handler) LoginRequest() {
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}
