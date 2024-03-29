package rpcx_test

import (
	"context"
	"fmt"
	"math"
	"net"
	"testing"
	"time"

	"github.com/vimcoders/go-driver/rpcx"

	"github.com/vimcoders/go-driver/pb"

	"github.com/vimcoders/go-driver/message"

	"github.com/vimcoders/go-driver/log"

	"github.com/vimcoders/go-driver/driver"

	"google.golang.org/protobuf/proto"
)

var total int64

type Handler struct {
	driver.Marshaler
	driver.Unmarshaler
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handler {
	encoder := message.NewProtobuf(message.GateMessages...)
	return &Handler{Marshaler: encoder, Unmarshaler: encoder}
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	connect := &rpcx.Connect{Conn: conn, Marshaler: x.Marshaler, Unmarshaler: x.Unmarshaler, Timeout: time.Second * 30}
	connect.OnMessage = func(request proto.Message) (proto.Message, error) {
		total++
		log.Info(total)
		return request, nil
	}
	go connect.Read(ctx)
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}

func TestListen(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp4", ":18888")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	go rpcx.ListenAndServe(context.Background(), listener, MakeHandler())
	conn, err := net.Dial("tcp", ":18888")
	if err != nil {
		fmt.Println(err)
		return
	}
	client := rpcx.NewClient(conn, message.GateMessages)
	go func() {
		for i := 0; i < math.MaxInt64; i++ {
			client.Call(context.Background(), 100, &pb.LoginRequest{Token: "token"})
		}
	}()
	go func() {
		for i := 0; i < math.MaxInt64; i++ {
			client.Call(context.Background(), 100, &pb.LoginRequest{Token: "token"})
		}
	}()
	go func() {
		for i := 0; i < math.MaxInt64; i++ {
			client.Call(context.Background(), 100, &pb.LoginRequest{Token: "token"})
		}
	}()
	go func() {
		for i := 0; i < math.MaxInt64; i++ {
			client.Call(context.Background(), 100, &pb.LoginRequest{Token: "token"})
		}
	}()
	go func() {
		for i := 0; i < math.MaxInt64; i++ {
			client.Call(context.Background(), 100, &pb.LoginRequest{Token: "token"})
		}
	}()
	go func() {
		for i := 0; i < math.MaxInt64; i++ {
			client.Call(context.Background(), 100, &pb.LoginRequest{Token: "token"})
		}
	}()
	go func() {
		for i := 0; i < math.MaxInt64; i++ {
			client.Call(context.Background(), 100, &pb.LoginRequest{Token: "token"})
		}
	}()
	for i := 0; i < math.MaxInt64; i++ {
		client.Call(context.Background(), 100, &pb.LoginRequest{Token: "token"})
	}
}
