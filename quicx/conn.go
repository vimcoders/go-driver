package quicx

import (
	"net"
	"syscall"
	"time"

	"github.com/quic-go/quic-go"
)

type Quic struct {
	conn   *net.UDPConn
	qconn  quic.Connection
	stream quic.Stream
}

// Read implements the Conn Read method.
func (x *Quic) Read(b []byte) (int, error) {
	return x.stream.Read(b)
}

// Write implements the Conn Write method.
func (x *Quic) Write(b []byte) (int, error) {
	return x.stream.Write(b)
}

// LocalAddr returns the local network address.
func (x *Quic) LocalAddr() net.Addr {
	return x.qconn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (x *Quic) RemoteAddr() net.Addr {
	return x.qconn.RemoteAddr()
}

// Close closes the connection.
func (x *Quic) Close() error {
	if x.stream != nil {
		return x.stream.Close()
	}
	return nil
}

// SetDeadline sets the deadline associated with the listener. A zero time value disables the deadline.
func (x *Quic) SetDeadline(t time.Time) error {
	return x.conn.SetDeadline(t)
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (x *Quic) SetReadDeadline(t time.Time) error {
	return x.conn.SetReadDeadline(t)
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (x *Quic) SetWriteDeadline(t time.Time) error {
	return x.conn.SetWriteDeadline(t)
}

// SetReadBuffer sets the size of the operating system's receive buffer associated with the connection.
func (x *Quic) SetReadBuffer(bytes int) error {
	return x.conn.SetReadBuffer(bytes)
}

// SetWriteBuffer sets the size of the operating system's transmit buffer associated with the connection.
func (x *Quic) SetWriteBuffer(bytes int) error {
	return x.conn.SetWriteBuffer(bytes)
}

// SyscallConn returns a raw network connection. This implements the syscall.Conn interface.
func (x *Quic) SyscallConn() (syscall.RawConn, error) {
	return x.conn.SyscallConn()
}

const (
	ReaderBuffsize = 8 * 1024
	WriterBuffsize = 8 * 1024
	Header         = 4
	Timeout        = time.Second * 16
)
