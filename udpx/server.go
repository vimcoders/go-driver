package udpx

import (
	"context"
	"net"
	"net/netip"

	"go-driver/driver"
	"go-driver/log"
)

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(ctx context.Context, addr string, handler driver.Handler) {
	addrPort, err := netip.ParseAddrPort(addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", net.UDPAddrFromAddrPort(addrPort))
	if err != nil {
		log.Error("ListenAndServe", err.Error())
	}
	log.Debug(conn.RemoteAddr())
	handler.Handle(ctx, conn)
}
