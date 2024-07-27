package handle

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

	"go-driver/app/logic/driver"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/mongox"
	"go-driver/pb"

	etcd "go.etcd.io/etcd/client/v3"
)

var handler = &Handle{}

type Handle struct {
	*mongox.Mongo
	*etcd.Client
	Users []*driver.User
	driver.Option
	total uint64
	unix  int64
	c     grpcx.Client
	sync.RWMutex
	pb.HandlerServer
}

// MakeHandler creates a Handler instance
func MakeHandler(opt driver.Option) *Handle {
	log.Info("etcd endpoints:", opt.Etcd.Endpoints)
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{opt.Etcd.Endpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	log.Info(opt.Mongo.Host, opt.Mongo.DB)
	mongo, err := mongox.Connect(opt.Mongo.Host, opt.Mongo.DB)
	if err != nil {
		panic(err)
	}
	handler.Mongo = mongo
	handler.Client = cli
	handler.Option = opt
	return handler
}

// Handle receives and executes redis commands
func (x *Handle) Handle(ctx context.Context, conn net.Conn) {
	log.Infof("new conn %s", conn.RemoteAddr().String())
	cli := grpcx.NewClient(conn, pb.Handler_ServiceDesc)
	if err := cli.Register(x); err != nil {
		log.Error(err.Error())
	}
	for i := 0; i < 1; i++ {
		go cli.Keeplive(context.Background(), &pb.PingRequest{})
	}
	x.unix = time.Now().Unix()
	x.c = cli
}

func (x *Handle) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
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
func (x *Handle) Close() error {
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
