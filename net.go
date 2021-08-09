package driver

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"time"
)

type Writer struct {
	net.Conn
	*bytes.Buffer
	t time.Duration
}

func (w *Writer) Write(p []byte) (n int, err error) {
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

func NewWriter(c net.Conn) io.Writer {
	return &Writer{c, bytes.NewBuffer(make([]byte, 256)), time.Second * 15}
}

type Reader struct {
	net.Conn
	*bufio.Reader
	t time.Duration
}

func (r *Reader) Read() (p []byte, err error) {
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

func NewReader(c net.Conn) *Reader {
	return &Reader{c, bufio.NewReaderSize(c, 512), time.Second * 5}
}
