package tcp

import (
	"go-driver/driver"
	"net"
)

func Dial(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, Option{Messages: driver.Messages}), nil
}
