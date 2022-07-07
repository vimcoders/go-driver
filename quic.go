package driver

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/lucas-clemente/quic-go"
)

type Quic struct {
	quic.Connection
	C chan []byte
}

func (c *Quic) Push() (err error) {
	for {
		select {
		case b, ok := <-c.C:
			if !ok {
				return errors.New("close")
			}
			stream, err := c.Connection.OpenStreamSync(context.Background())
			if err != nil {
				return err
			}
			_, err = stream.Write(b)
			if err != nil {
				return err
			}
			stream.Close()
		}
	}
	return nil
}

type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	return w.Writer.Write(b)
}

func (c *Quic) Pull() (err error) {
	for {
		stream, err := c.Connection.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(loggingWriter{stream}, stream)
		stream.Close()
	}
	return nil
}

func (c *Quic) Close() (err error) {
	close(c.C)
	return nil
}
