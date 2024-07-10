package grpcx

import (
	"crypto/tls"
	"fmt"
	"go-driver/quicx"
	"net"
	"time"

	"google.golang.org/grpc"
)

func Dial(network string, addr string, desc grpc.ServiceDesc) (Client, error) {
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
		return newClient(conn, desc), nil
	case "tcp":
		fallthrough
	case "tcp4":
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		return newClient(conn, desc), nil
	}
	return nil, fmt.Errorf("%s unkonw", network)
}
