package tcp

import (
	"net"
	"time"
)

func Dial(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, Option{Timeout: time.Minute, Buffsize: 1024}), nil
}
