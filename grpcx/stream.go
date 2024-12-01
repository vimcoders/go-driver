package grpcx

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/vimcoders/go-driver/log"

	"google.golang.org/protobuf/proto"
)

type stream struct {
	net.Conn
	seq     uint32
	timeout time.Duration
	signal  chan Message
}

func (x *stream) invoke(iMessage Message) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	x.signal <- iMessage
	return nil
}

func (x *stream) push(ctx context.Context, method uint16, req proto.Message) (Message, error) {
	buf, err := encode(x.seq, method, req)
	if err != nil {
		return nil, err
	}
	if err := x.SetWriteDeadline(time.Now().Add(x.timeout)); err != nil {
		return nil, err
	}
	if _, err := buf.WriteTo(x.Conn); err != nil {
		return nil, err
	}
	select {
	case iMessage := <-x.signal:
		return iMessage, nil
	case <-ctx.Done():
		close(x.signal)
		return nil, errors.New("timeout")
	}
}
