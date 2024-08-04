package tcp_test

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"testing"

	"go-driver/driver"
	"go-driver/pb"
	"go-driver/tcp"

	"go-driver/log"

	"google.golang.org/protobuf/proto"
)

type Buffer = driver.Buffer

func NewBuffer(size int) Buffer {
	return driver.NewBuffer(size)
}

type Handler struct {
	messages []proto.Message
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handler {
	return &Handler{
		messages: driver.Messages,
	}
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	cli := tcp.NewClient(conn, tcp.Option{Messages: x.messages})
	cli.Register(x)
}

func (x *Handler) ServeTCP(ctx context.Context, req proto.Message) error {
	log.Debug(req, "ServeTCP")
	return nil
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}

func TestMain(m *testing.M) {
	addr, err := net.ResolveTCPAddr("tcp4", ":18888")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	handler := MakeHandler()
	for i := 0; i < runtime.NumCPU(); i++ {
		go tcp.ListenAndServe(context.Background(), listener, handler)
	}
	m.Run()
}

func BenchmarkTCP(b *testing.B) {
	conn, err := net.Dial("tcp", ":18888")
	if err != nil {
		fmt.Println(err)
		return
	}
	cli := tcp.NewClient(conn, tcp.Option{Messages: driver.Messages})
	cli.Register(MakeHandler())
	for i := 0; i < b.N; i++ {
		cli.Go(context.Background(), &pb.LoginRequest{Token: "token"})
	}
}
