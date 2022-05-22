package driver

import (
	"errors"
	"net"
	"time"
)

type Session struct {
	net.Conn
	OnMessage func(b []byte) error
	C         chan []byte
}

func (s *Session) Push() (err error) {
	buffer := make([]byte, 512)
	for {
		select {
		case b, ok := <-s.C:
			if !ok {
				return errors.New("close")
			}
			length := len(b) + 2
			buffer[0] = uint8(length >> 8)
			buffer[1] = uint8(length)
			copy(buffer[2:], b)
			if err = s.SetWriteDeadline(Timeout()); err != nil {
				return err
			}
			if _, err = s.Conn.Write(buffer[:length]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Session) Pull() (err error) {
	buffer := make([]byte, 512)
	length := 0
	for {
		if err = s.SetReadDeadline(Timeout()); err != nil {
			return err
		}
		n, err := s.Read(buffer[length:])
		if err != nil {
			return err
		}
		length += n
		header := int(int(buffer[0])<<8 | int(buffer[1]))
		if header > 512 {
			return errors.New("header too long")
		}
		if length < header {
			continue
		}
		if err := s.OnMessage(buffer[2:header]); err != nil {
			return err
		}
		copy(buffer, buffer[header:])
		length -= header
	}

	return nil
}

func (s *Session) Close() (err error) {
	if err = s.Conn.Close(); err != nil {
		return err
	}
	close(s.C)

	return nil
}

func Timeout() time.Time {
	return time.Now().Add(15 * time.Second)
}
