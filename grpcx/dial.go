package grpcx

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/vimcoders/go-driver/quicx"
)

func Dial(network string, addr string, opt Option) (Client, error) {
	switch network {
	case "udp":
		conn, err := quicx.Dial(addr, &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"quic-echo-example"},
			MaxVersion:         tls.VersionTLS13,
		}, &quicx.Config{
			MaxIdleTimeout: time.Minute,
		})
		if err != nil {
			return nil, err
		}
		return newClient(conn, opt), nil
	case "tcp":
		fallthrough
	case "tcp4":
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		return newClient(conn, opt), nil
	}
	return nil, fmt.Errorf("%s unkonw", network)
}
