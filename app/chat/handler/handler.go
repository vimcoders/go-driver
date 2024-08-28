package handler

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"go-driver/app/chat/driver"
	"go-driver/etcdx"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/quicx"
	"math/big"
	"net"
	"os"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

type Option = driver.Option

type Handler struct {
	Option
	*etcd.Client
	pb.ChatServer
}

func MakeHandler(ctx context.Context) *Handler {
	h := &Handler{}
	if err := h.parse(); err != nil {
		panic(err)
	}
	if err := h.connect(ctx); err != nil {
		panic(err)
	}
	return h
}

func (x *Handler) parse() error {
	var fileName string
	flag.StringVar(&fileName, "option", "chat.conf", "chat.conf")
	flag.Parse()
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(ymalBytes, &x.Option); err != nil {
		return err
	}
	return nil
}

func (x *Handler) connect(_ context.Context) error {
	log.Info(x.Etcd.Endpoints)
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{x.Etcd.Endpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	x.Client = cli
	return nil
}

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
	b, err := json.Marshal(&etcdx.Service{Internet: x.QUIC.Internet, Local: x.QUIC.Local})
	if err != nil {
		panic(err)
	}
	if _, err := x.Client.Put(ctx, x.Etcd.Join("chat"), string(b)); err != nil {
		panic(err)
	}
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
