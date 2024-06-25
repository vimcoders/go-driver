package rpcx_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"testing"
	"time"

	"go-driver/log"
	"go-driver/pb"
	"go-driver/quicx"
	"go-driver/rpcx"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handler {
	return &Handler{}
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	handle := &rpcx.Handle{
		Conn:     conn,
		Buffsize: 16 * 1024,
		Timeout:  time.Second * 120,
	}
	handle.Register(x)
}

func (x *Handler) PingRequest(ctx context.Context, request *pb.PingRequest) (proto.Message, error) {
	return &pb.PingResponse{Message: []byte("ping .")}, nil
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}

func TestMain(m *testing.M) {
	rect := new(Rect)
	rpc.Register(rect)
	rpc.HandleHTTP()
	go http.ListenAndServe(":9900", nil)
	addr, err := net.ResolveTCPAddr("tcp4", ":18888")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	go rpcx.ListenAndServe(context.Background(), listener, MakeHandler())
	listenerQuic, err := quicx.Listen("udp", ":8973", generateTLSConfig(), &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	go rpcx.ListenAndServe(context.Background(), listenerQuic, MakeHandler())
	s := server.NewServer()
	s.RegisterName("Arith", new(Arith), "")
	go s.Serve("tcp", ":9972")
	m.Run()
}

func TestTcp(t *testing.T) {
	conn, err := net.Dial("tcp", ":18888")
	if err != nil {
		fmt.Println(err)
		return
	}
	client := rpcx.NewClient(conn)
	for i := 0; i < 5000; i++ {
		go func() {
			for i := 0; i < math.MaxInt64; i++ {
				var reply pb.PingResponse
				client.Call(context.Background(), &pb.PingRequest{Message: []byte("ping")}, &reply)
			}
		}()
	}
}

func TestQuic(t *testing.T) {
	conn, err := quicx.Dial("127.0.0.1:8972", &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
		MaxVersion:         tls.VersionTLS13,
	}, &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	client := rpcx.NewClient(conn)
	for i := 0; i < 5000; i++ {
		go func() {
			for i := 0; i < math.MaxInt32; i++ {
				var reply pb.PingResponse
				client.Call(context.Background(), &pb.PingRequest{Message: []byte("ping")}, &reply)
			}
		}()
	}
}

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith struct {
	Total int64
}

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	t.Total++
	// if t.Total%200000 == 0 {
	// 	//log.Debug(t.Total)
	// }
	return nil
}

func (t *Arith) Add(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Say(ctx context.Context, args *string, reply *string) error {
	*reply = "hello " + *args
	return nil
}

func TestRPCX(t *testing.T) {
	d, err := client.NewPeer2PeerDiscovery("tcp@"+":9972", "")
	if err != nil {
		t.Error(err.Error())
		return
	}
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	args := &Args{
		A: 10,
		B: 20,
	}
	for i := 0; i < 5000; i++ {
		go func() {
			for i := 0; i < math.MaxInt32; i++ {
				reply := &Reply{}
				xclient.Call(context.Background(), "Mul", args, reply)
			}
		}()
	}
	c := time.NewTimer(time.Minute)
	defer c.Stop()
	<-c.C
}

func generateTLSConfig() *tls.Config {
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

type Params struct {
	Width, Height int
}

type Rect struct {
	Total int64
}

// RPC服务端方法，求矩形面积
func (r *Rect) Area(p Params, ret *int) error {
	*ret = p.Height * p.Width
	r.Total++
	return nil
}

// 周长
func (r *Rect) Perimeter(p Params, ret *int) error {
	*ret = (p.Height + p.Width) * 2
	return nil
}

func TestRPC(t *testing.T) {
	conn, err := rpc.DialHTTP("tcp", ":8000")
	if err != nil {
		log.Error(err.Error())
	}
	for i := 0; i < 5000; i++ {
		go func() {
			for i := 0; i < math.MaxInt32; i++ {
				ret := 0
				if err := conn.Call("Rect.Area", Params{50, 100}, &ret); err != nil {
					log.Error(err.Error())
				}
			}
		}()
	}
	c := time.NewTimer(time.Minute)
	defer c.Stop()
	<-c.C
}

func BenchmarkRPC(b *testing.B) {
	conn, err := rpc.DialHTTP("tcp", ":9900")
	if err != nil {
		b.Error(err.Error())
		return
	}
	defer conn.Close()
	for i := 0; i < b.N; i++ {
		ret := 0
		if err := conn.Call("Rect.Area", Params{50, 100}, &ret); err != nil {
			b.Error(err.Error())
			return
		}
	}
}

func BenchmarkTCP(b *testing.B) {
	b.Log("Benchmark_RPC")
	conn, err := net.Dial("tcp", ":18888")
	if err != nil {
		fmt.Println(err)
		return
	}
	client := rpcx.NewClient(conn)
	for i := 0; i < b.N; i++ {
		var reply pb.PingResponse
		client.Call(context.Background(), &pb.PingRequest{Message: []byte("ping")}, &reply)
	}
}

func BenchmarkQuic(b *testing.B) {
	b.Log("Benchmark_RPC")
	conn, err := quicx.Dial("127.0.0.1:8973", &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
		MaxVersion:         tls.VersionTLS13,
	}, &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	client := rpcx.NewClient(conn)
	for i := 0; i < b.N; i++ {
		var reply pb.PingResponse
		client.Call(context.Background(), &pb.PingRequest{Message: []byte("ping")}, &reply)
	}
}
