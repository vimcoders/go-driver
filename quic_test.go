package driver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/lucas-clemente/quic-go"
)

func init_quic() {
	listener, err := quic.ListenAddr("localhost:9999", generateTLSConfig(), nil)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			c, err := listener.Accept(context.Background())
			if err != nil {
				panic(err)
			}
			NewQuic(c)
		}
	}()
}

func TestQuic(t *testing.T) {
	init_quic()
	var waitGroup sync.WaitGroup
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	for i := 0; i < 1; i++ {
		waitGroup.Add(1)
		c, err := quic.DialAddr("localhost:9999", tlsConf, nil)
		if err != nil {
			continue
		}
		quic := NewQuic(c)
		go func() {
			ticker := time.NewTicker(time.Second)
			for range ticker.C {
				quic.C <- []byte("hello")
			}
		}()
	}
	waitGroup.Wait()
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
	}
}
