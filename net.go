package driver

import (
	"errors"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	OnMessage func(b []byte) error
	C         chan []byte
}

func (c *Conn) Push() (err error) {
	buffer := make([]byte, 512)
	for {
		select {
		case b, ok := <-c.C:
			if !ok {
				return errors.New("close")
			}
			length := len(b) + 2
			buffer[0] = uint8(length >> 8)
			buffer[1] = uint8(length)
			copy(buffer[2:], b)
			if err = c.SetWriteDeadline(Timeout()); err != nil {
				return err
			}
			if _, err = c.Write(buffer[:length]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Conn) Pull() (err error) {
	buffer := make([]byte, 512)
	length := 0
	for {
		if err = c.SetReadDeadline(Timeout()); err != nil {
			return err
		}
		n, err := c.Read(buffer[length:])
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
		if err := c.OnMessage(buffer[2:header]); err != nil {
			return err
		}
		copy(buffer, buffer[header:])
		length -= header
	}

	return nil
}

func (c *Conn) Close() (err error) {
	if err = c.Conn.Close(); err != nil {
		return err
	}
	close(c.C)
	return nil
}

func Timeout() time.Time {
	return time.Now().Add(15 * time.Second)
}
