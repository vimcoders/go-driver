package kcpx

import (
	"net"

	"github.com/xtaci/kcp-go"
)

// Listen creates a QUIC listener on the given network interface
func Listen(laddr string) (net.Listener, error) {
	return kcp.Listen(laddr)
}
