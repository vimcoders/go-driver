package handler

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"go-driver/udpx"
	"math/big"
)

func (x *Handler) ListenAndServe(ctx context.Context) {
	//listener, err := tls.Listen("tcp", x.TCP.LocalAddr, generateTLSConfig())
	//listener, err := net.Listen("tcp", x.TCP.LocalAddr)
	// listener, err := quicx.Listen("udp", x.QUIC.LocalAddr, generateTLSConfig(), &quicx.Config{
	// 	MaxIdleTimeout: time.Minute,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	//for i := 0; i < runtime.NumCPU(); i++ {
	//go quicx.ListenAndServe(ctx, listener, x)
	go udpx.ListenAndServe(ctx, x.TCP.Internet, x)
	//}
	//log.Infof("running %s", listener.Addr().String())
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
		//MaxVersion:   tls.VersionTLS12,
	}
}
