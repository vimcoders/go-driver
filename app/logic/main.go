// 逻辑服务 处理请求
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"go-driver/app/logic/driver"
	"go-driver/app/logic/handle"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/quicx"
)

func main() {
	log.Info("runtime.NumCPU: ", runtime.NumCPU())
	option := driver.ReadOption()
	log.Info("Unmarshal YAML DOWN")
	handler := handle.MakeHandler(option)
	defer handler.Close()
	// addr, err := net.ResolveTCPAddr("tcp4", opt.Addr.Port)
	// if err != nil {
	// 	panic(err)
	// }
	// listener, err := net.ListenTCP("tcp", addr)
	listener, err := quicx.Listen("udp", option.QUIC.Port, GenerateTLSConfig(), &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go grpcx.ListenAndServe(ctx, listener, handler)
	log.Infof("RUNNING %s", listener.Addr().String())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-quit
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
			log.Info("shutdown ->", s.String())
			cancel()
		default:
			log.Info("os.Signal ->", s.String())
			continue
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
