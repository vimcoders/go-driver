package udpx

import (
	"net"
	"net/netip"
)

func Dial(network string, addr string, opt Option) (Client, error) {
	addrPort, err := netip.ParseAddrPort(addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP(network, &net.UDPAddr{}, net.UDPAddrFromAddrPort(addrPort))
	if err != nil {
		return nil, err
	}
	return newClient(conn, opt), nil
}
