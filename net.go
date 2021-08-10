package driver

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"time"
)

type writer struct {
	net.Conn
	*bytes.Buffer
	t time.Duration
}

func (w *writer) Write(p []byte) (n int, err error) {
	defer w.Reset()

	length := len(p)

	var header [2]byte

	header[0] = uint8(length >> 8)
	header[1] = uint8(length)

	if err := w.SetWriteDeadline(time.Now().Add(w.t)); err != nil {
		return 0, err
	}

	if n, err := w.Buffer.Write(header[:]); err != nil {
		return n, err
	}

	if n, err := w.Buffer.Write(p); err != nil {
		return n, err
	}

	return w.Conn.Write(w.Buffer.Bytes())
}

func NewWriter(c net.Conn, b *bytes.Buffer, t time.Duration) io.Writer {
	return &writer{c, b, t}
}

type reader struct {
	net.Conn
	*bufio.Reader
	t time.Duration
}

func (r *reader) Read() (p []byte, err error) {
	header := make([]byte, 2)

	if err := r.SetReadDeadline(time.Now().Add(r.t)); err != nil {
		return nil, err
	}

	if _, err := r.Reader.Read(header); err != nil {
		return nil, err
	}

	length := uint16(uint16(header[0])<<8 | uint16(header[1]))

	if err := r.SetReadDeadline(time.Now().Add(r.t)); err != nil {
		return nil, err
	}

	return r.Reader.Peek(int(length))
}

func NewReader(c net.Conn, b *bufio.Reader, t time.Duration) Reader {
	return &reader{c, b, t}
}
