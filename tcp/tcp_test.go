package tcp_test

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"go-driver/driver"
	"go-driver/log"
	"go-driver/tcp"
)

type Buffer []byte

func NewBuffer(size int) Buffer {
	return make([]byte, size)
}

func (b *Buffer) Reset() {
	*b = (*b)[:0]
}

func (b *Buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *Buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

type TcpSession struct {
	Timeout        time.Duration
	Header         int
	ReaderBuffsize int
	net.Conn
}

func (x *TcpSession) Pull(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
		if err != nil {
			fmt.Println(err.Error())
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
		x.ServeTCP(x.Conn, message[x.Header:])
		if _, err := buffer.Discard(len(message)); err != nil {
			return err
		}
	}
}

func (x *TcpSession) ServeTCP(w io.Writer, request []byte) {
	w.Write(request)
}

func (x *TcpSession) Close() error {
	return nil
}

type Handler struct {
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handler {
	return &Handler{}
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	newSession := &TcpSession{
		ReaderBuffsize: 12 * 1026,
		Header:         4,
		Timeout:        time.Minute,
		Conn:           conn,
	}
	go newSession.Pull(context.Background())
}

func (x *Handler) ServeTCP(w io.Writer, request []byte) {
	w.Write(request)
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}

type Session struct {
	total          int64
	ReaderBuffsize int
	Header         int
	Timeout        time.Duration
	net.Conn
}

func (x *Session) Pull(ctx context.Context) (err error) {
	defer func() {
		// if e := recover(); e != nil {
		// 	log.Error(fmt.Sprintf("%s", e))
		// 	debug.PrintStack()
		// }
		// if err != nil {
		// 	log.Error(err.Error(), x.RemoteAddr().String())
		// 	debug.PrintStack()
		// }
		// x.Close()
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
		x.ServeTCP(x.Conn, message[x.Header:])
		if _, err := buffer.Discard(len(message)); err != nil {
			return err
		}
	}
}

func (x *Session) ServeTCP(w io.Writer, request []byte) {
	w.Write(request)
}

func (x *Session) Close() error {
	return nil
}

func TestMain(m *testing.M) {
	addr, err := net.ResolveTCPAddr("tcp4", ":18888")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	handler := MakeHandler()
	for i := 0; i < runtime.NumCPU(); i++ {
		go tcp.ListenAndServe(context.Background(), listener, handler)
	}
	m.Run()
}

func BenchmarkTCP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go func() {
			conn, err := net.Dial("tcp", ":18888")
			if err != nil {
				fmt.Println(err)
				return
			}
			for k := 0; k < b.N; k++ {
				b := []byte("12341342341241")
				buffer := make(driver.Buffer, 4)
				binary.BigEndian.PutUint32(buffer, uint32(len(b)))
				buffer.Write(b)
				conn.Write(buffer)
				conn.Read(buffer)
			}
		}()
	}
}

// func generateTLSConfig() *tls.Config {
// 	key, err := rsa.GenerateKey(rand.Reader, 1024)
// 	if err != nil {
// 		panic(err)
// 	}
// 	template := x509.Certificate{SerialNumber: big.NewInt(1)}
// 	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
// 	if err != nil {
// 		panic(err)
// 	}
// 	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
// 	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

// 	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &tls.Config{
// 		Certificates: []tls.Certificate{tlsCert},
// 		//NextProtos:   []string{"http/1.1"},
// 		MaxVersion: tls.VersionTLS13,
// 	}
// }
