package rpcx_test

import (
	"context"
	"fmt"
	"math"
	"net"
	"testing"
	"time"

	"github.com/vimcoders/go-driver/driver"
	"github.com/vimcoders/go-driver/log"
	"github.com/vimcoders/go-driver/message"
	"github.com/vimcoders/go-driver/pb"
	"github.com/vimcoders/go-driver/rpcx"
	"google.golang.org/protobuf/proto"
)

var total int64

type Handler struct {
	driver.Marshaler
	driver.Unmarshaler
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handler {
	encoder := message.NewProtobuf(message.Messages...)
	return &Handler{Marshaler: encoder, Unmarshaler: encoder}
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	connect := &rpcx.Connect{Conn: conn, Timeout: time.Second * 30}
	connect.OnMessage = func(w rpcx.ResponseWriter, request *rpcx.Request) {
		total++
		log.Info(total)
		w.Write(request.Message)
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
	client := rpcx.NewClient(conn, message.Messages)
	go func() {
		b, err := proto.Marshal(&pb.LoginRequest{Token: "login"})
		if err != nil {
			log.Error(err.Error())
		}
		for i := 0; i < math.MaxInt64; i++ {
			client.Call(context.Background(), &rpcx.Request{Message: b})
		}
	}()
	b, err := proto.Marshal(&pb.LoginRequest{Token: "login"})
	if err != nil {
		log.Error(err.Error())
	}
	for i := 0; i < math.MaxInt64; i++ {
		client.Call(context.Background(), &rpcx.Request{Message: b})
	}
}
