// 逻辑服务 处理请求
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"math/big"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"go-driver/app/logic/driver"
	"go-driver/app/logic/handle"
	"go-driver/etcdx"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/quicx"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

func main() {
	log.Info("runtime.NumCPU: ", runtime.NumCPU())
	var fileName string
	flag.StringVar(&fileName, "conf", "./logic.conf", "logic.conf")
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	log.Info("READ YAML DOWN")
	var opt driver.YAML
	if err := yaml.Unmarshal(ymalBytes, &opt); err != nil {
		panic(err.Error())
	}
	log.Info("Unmarshal YAML DOWN")
	handler := handle.MakeHandler(&opt)
	defer handler.Close()
	// addr, err := net.ResolveTCPAddr("tcp4", opt.Addr.Port)
	// if err != nil {
	// 	panic(err)
	// }
	// listener, err := net.ListenTCP("tcp", addr)
	listener, err := quicx.Listen("udp", opt.QUIC.Port, GenerateTLSConfig(), &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	service := etcdx.Service{
		Kind: "QUIC",
		Addr: opt.QUIC.Internet,
	}
	b, err := json.Marshal(service)
	if err != nil {
		panic(err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	key := opt.Etcd.Version + "/service/logic/" + uuid.NewString()
	if _, err := handler.Client.Put(ctx, key, string(b)); err != nil {
		panic(err.Error())
	}
	log.Info("ETCD Put down")
	go grpcx.ListenAndServe(ctx, listener, handler)
	log.Infof("RUNNING %s %s", listener.Addr().String(), key)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	log.Info("SHUTDOWN", <-quit)
	handler.Client.Delete(ctx, key)
	log.Info("ETCD DEL DOWN..")
	cancel()
	log.Info("cancel() DOWN....")
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
