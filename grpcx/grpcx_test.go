// goos: linux
// goarch: amd64
// pkg: go-driver/grpcx
// cpu: 12th Gen Intel(R) Core(TM) i3-12100F
// BenchmarkTCP-8
// 468010  request/s  NumGoroutine  9
//
//	667417              1685 ns/op             248 B/op          5 allocs/op
//
// BenchmarkQUIC-8
// 964540  request/s  NumGoroutine  36
// 1518244  request/s  NumGoroutine  43
//
//	1831996               657.2 ns/op           281 B/op          5 allocs/op
//
// PASS

package grpcx_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go-driver/grpcx"
	"go-driver/pb"
	"go-driver/quicx"
	"math/big"
	"net"
	"net/http"
	"runtime"
	"testing"
	"time"
)

type Handle struct {
	total uint64
	unix  int64
	pb.HandlerServer
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handle {
	return &Handle{}
}

// Handle receives and executes redis commands
func (x *Handle) Handle(ctx context.Context, conn net.Conn) {
	cli := grpcx.NewClient(conn, grpcx.Option{ServiceDesc: pb.Handler_ServiceDesc})
	//go cli.Keeplive(ctx, &pb.PingRequest{})
	cli.Register(x)
}

func (x *Handle) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	unix := time.Now().Unix()
	x.total++
	if unix != x.unix {
		fmt.Println(x.total, " request/s", " NumGoroutine ", runtime.NumGoroutine())
		x.total = 0
		x.unix = unix
	}
	return &pb.PingResponse{Message: req.Message}, nil
}

func (x *Handle) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{Code: http.StatusOK}, nil
}

// Close stops handler
func (x *Handle) Close() error {
	return nil
}

func TestMain(m *testing.M) {
	addr, err := net.ResolveTCPAddr("tcp4", ":28888")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	go grpcx.ListenAndServe(context.Background(), listener, MakeHandler())
	qlistener, err := quicx.Listen("udp", ":28889", GenerateTLSConfig(), &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	go grpcx.ListenAndServe(context.Background(), qlistener, MakeHandler())
	m.Run()
}

func BenchmarkTCP(b *testing.B) {
	fmt.Println(runtime.NumCPU())
	cli, err := grpcx.Dial("tcp", "127.0.0.1:28888", grpcx.Option{ServiceDesc: pb.Handler_ServiceDesc})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cli.Register(MakeHandler())
	for i := 0; i < b.N; i++ {
		if err := cli.Go(context.Background(), "Ping", &pb.PingRequest{}); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func BenchmarkQUIC(b *testing.B) {
	fmt.Println(runtime.NumCPU())
	cli, err := grpcx.Dial("udp", "127.0.0.1:28889", grpcx.Option{ServiceDesc: pb.Handler_ServiceDesc})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cli.Register(MakeHandler())
	for i := 0; i < b.N; i++ {
		if err := cli.Go(context.Background(), "Ping", &pb.PingRequest{}); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func BenchmarkMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var message grpcx.Message
		message.WriteUint16(1)
		message.WriteUint32(1)
		message.WriteUint16(1)
		message.Write([]byte{1, 2, 3, 4, 5, 6, 7})
	}
}

func BenchmarkPing(b *testing.B) {
	fmt.Println(runtime.NumCPU())
	cli, err := grpcx.Dial("udp", "127.0.0.1:28889", grpcx.Option{ServiceDesc: pb.Handler_ServiceDesc})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cli.Register(MakeHandler())
	client := pb.NewHandlerClient(cli)
	for i := 0; i < b.N; i++ {
		if _, err := client.Ping(context.Background(), &pb.PingRequest{}); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func GenerateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
		MaxVersion:   tls.VersionTLS13,
	}
}
