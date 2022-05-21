package driver

import (
	"bufio"
	"bytes"
	"net"
	"time"
)

type session struct {
	net.Conn
	*bytes.Buffer
	*bufio.Reader
}

func (s *session) Write(b []byte) (n int, err error) {
	length := len(b)

	var header [2]byte

	header[0] = uint8(length >> 8)
	header[1] = uint8(length)

	if _, err := s.Buffer.Write(header[:]); err != nil {
		return 0, err
	}

	if _, err := s.Buffer.Write(b); err != nil {
		return 0, err
	}

	return s.Conn.Write(s.Buffer.Bytes())
}

func (s *session) Read() (b []byte, err error) {
	header := make([]byte, 2)

	if _, err := s.Reader.Read(header); err != nil {
		return nil, err
	}

	length := uint16(uint16(header[0])<<8 | uint16(header[1]))

	if err := s.SetReadDeadline(time.Now().Add(time.Duration(10))); err != nil {
		return nil, err
	}

	return s.Reader.Peek(int(length))
}
