package handler

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"
	"runtime"
	"sync"
	"time"

	"go-driver/app/parkour/driver"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/mongox"
	"go-driver/pb"

	etcd "go.etcd.io/etcd/client/v3"
)

var handler *Handler

type Handler struct {
	*mongox.Mongo
	*etcd.Client
	Users []*driver.User
	Option
	total uint64
	unix  int64
	c     grpcx.Client
	sync.RWMutex
	pb.ParkourServer
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
	return h
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	log.Infof("new conn %s", conn.RemoteAddr().String())
	cli := grpcx.NewClient(conn, grpcx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	if err := cli.Register(ctx, x); err != nil {
		log.Error(err.Error())
	}
	x.unix = time.Now().Unix()
	x.c = cli
}

func (x *Handler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	unix := time.Now().Unix()
	x.total++
	if unix != x.unix {
		log.Debug(x.total, " request/s", " NumGoroutine ", runtime.NumGoroutine())
		x.total = 0
		x.unix = unix
	}
	return &pb.PingResponse{Message: req.Message}, nil
}

// Close stops handler
func (x *Handler) Close() error {
	for i := 0; i < len(x.Users); i++ {
		// x.Mongo.Insert(x.Users[i])
		// x.Mongo.Update(x.Users[i])
	}
	return nil
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
