package rpcx

import (
	"net"
)

type Quic struct {
	net.Conn
}
