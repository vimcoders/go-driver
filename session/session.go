package session

import (
	"context"
	"errors"
	"net"
	"time"
)

type Session struct {
	Id int64
	net.Conn
	OnMessage func(b []byte) error
	d         map[interface{}]interface{}
	C         chan []byte
	context.CancelFunc
}

func (s *Session) Set(key, value interface{}) error {
	s.d[key] = value
	return nil
}

func (s *Session) Get(key interface{}) interface{} {
	if v, ok := s.d[key]; ok {
		return v
	}
	return nil
}

func (s *Session) Delete(key interface{}) error {
	delete(s.d, key)
	return nil
}

func (s *Session) Push(ctx context.Context) (err error) {
	defer s.Close()
	buffer := make([]byte, 512)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
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

func (s *Session) Pull(ctx context.Context) (err error) {
	defer s.Close()
	buffer := make([]byte, 512)
	length := 0
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
			header := int(int(buffer[0])<<8 | int(buffer[1]))
			if header > 512 {
				return errors.New("header too long")
			}
			for header <= 0 || length < header {
				if err = s.SetReadDeadline(Timeout()); err != nil {
					return err
				}
				n, err := s.Read(buffer[length:])
				if err != nil {
					return err
				}
				length += n
				header = int(int(buffer[0])<<8 | int(buffer[1]))
				if header > 512 {
					return errors.New("header too long")
				}
			}
			if err := s.OnMessage(buffer[2:header]); err != nil {
				return err
			}
			copy(buffer, buffer[header:])
			length -= header
		}
	}
}

func (s *Session) Close() (err error) {
	if err = s.Conn.Close(); err != nil {
		return err
	}
	s.CancelFunc()
	close(s.C)
	return nil
}

func Timeout() time.Time {
	return time.Now().Add(15 * time.Second)
}
