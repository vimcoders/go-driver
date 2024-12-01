package tcpx

import (
	"context"
	"runtime/debug"

	"github.com/vimcoders/go-driver/tcpx"

	"github.com/vimcoders/go-driver/pb"
)

type Client struct {
	tcpx.Client
}

func Dial(network string, addr string) (*Client, error) {
	client, err := tcpx.Dial(network, addr, tcpx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	if err != nil {
		return nil, err
	}
	c := &Client{Client: client}
	c.Register(context.Background(), nil)
	return c, nil
}

func (x *Client) BenchmarkTCP() (err error) {
	defer func() {
		debug.PrintStack()
	}()
	for {
		if err := x.Go(context.Background(), &pb.PingRequest{}); err != nil {
			panic(err)
		}
	}
}

func (x *Client) Handle(ctx context.Context, iMessage tcpx.Message) error {
	return nil
}
