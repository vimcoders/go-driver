package quicx_test

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"runtime/debug"
	"testing"
	"time"

	"github.com/vimcoders/go-driver/log"
	"github.com/vimcoders/go-driver/quicx"
)

type Handler struct {
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handler {
	return &Handler{}
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
}

func (x *Handler) ServeTCP(w io.Writer, request []byte) {
	w.Write(request)
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}

type QuicSession struct {
	ReaderBuffsize int
	Header         int
	Timeout        time.Duration
	net.Conn
}

func (x *QuicSession) Pull(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
		if err != nil {
			//log.Error(err.Error(), x.RemoteAddr().String())
			debug.PrintStack()
		}
		x.Close()
	}()
	buffer := bufio.NewReaderSize(x.Conn, x.ReaderBuffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		headerBytes, err := buffer.Peek(x.Header)
		if err != nil {
			return err
		}
		header := binary.BigEndian.Uint32(headerBytes)
		if int(header) > buffer.Size() {
			return fmt.Errorf("header %v too long", header)
		}
		message, err := buffer.Peek(int(header) + len(headerBytes))
		if err != nil {
			return err
		}
		if len(message) < x.Header {
			return errors.New("len(bodyBytes) < Header+Proto")
		}
		x.ServeQuic(x.Conn, message[x.Header:])
		if _, err := buffer.Discard(len(message)); err != nil {
			return err
		}
	}
}

func (x *QuicSession) ServeQuic(w io.Writer, request []byte) {
	w.Write(request)
}

func (x *QuicSession) Close() error {
	return nil
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
		//NextProtos:   []string{"http/1.1"},
		MaxVersion: tls.VersionTLS13,
	}
}

func TestMain(m *testing.M) {
	listenerQuic, err := quicx.Listen("udp", ":8972", generateTLSConfig(), &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	go quicx.ListenAndServe(context.Background(), listenerQuic, MakeHandler())
	m.Run()
}

func BenchmarkQuic(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		conn.Write([]byte("1234"))
	}
}
