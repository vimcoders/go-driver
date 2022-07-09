package driver

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucas-clemente/quic-go"
)

type Quic struct {
	quic.Connection
	OnMessage func(b []byte) error
	C         chan []byte
}

func NewQuic(c quic.Connection) *Quic {
	q := &Quic{
		Connection: c,
		OnMessage: func(b []byte) error {
			fmt.Println(string(b))
			return nil
		},
		C: make(chan []byte, 1),
	}
	go q.Pull()
	go q.Push()
	return q
}

func (c *Quic) Push() (err error) {
	stream, err := c.Connection.OpenStreamSync(context.Background())
	defer stream.Close()
	if err != nil {
		return err
	}
	for {
		select {
		case b, ok := <-c.C:
			if !ok {
				return errors.New("close")
			}
			_, err = stream.Write(b)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Quic) Pull() (err error) {
	r, err := c.Connection.AcceptStream(context.Background())
	defer r.Close()
	if err != nil {
		return err
	}
	b := make([]byte, 512)
	for {
		nr, err := r.Read(b)
		if err != nil {
			return err
		}
		if nr <= 0 {
			continue
		}
		if err := c.OnMessage(b); err != nil {
			return err
		}
	}
	return nil
}

func (c *Quic) Close() (err error) {
	close(c.C)
	return nil
}
