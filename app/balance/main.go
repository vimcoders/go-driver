// 负载均衡服务 请求转发到服务去处理
// 同时我们也会主动推送很多数据，它的推送速度就像彗星一样的快
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"math/big"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go-driver/app/balance/handle"
	"go-driver/conf"
	"go-driver/log"
	"go-driver/quicx"

	"gopkg.in/yaml.v3"
)

func main() {
	runtime.GOMAXPROCS(3)
	var fileName string
	flag.StringVar(&fileName, "conf", "./comet.conf", "comet.conf")
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	var opt conf.Conf
	if err := yaml.Unmarshal(ymalBytes, &opt); err != nil {
		panic(err.Error())
	}
	handler := handle.MakeHandler(opt)
	addr, err := net.ResolveTCPAddr("tcp4", opt.Addr.Port)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	// tcpAddr := listener.Addr().(*net.TCPAddr)
	// listener, err := quicx.Listen("udp", opt.Addr.Port, GenerateTLSConfig(), &quicx.Config{
	// 	MaxIdleTimeout: time.Minute,
	// })
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < runtime.NumCPU(); i++ {
		go quicx.ListenAndServe(ctx, listener, handler)
	}
	log.Infof("running %s", listener.Addr().String())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	log.Info("SHUTDOWN", <-quit)
	handler.Close()
	cancel()
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
