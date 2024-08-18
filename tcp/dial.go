package tcp

import (
	"go-driver/driver"
	"net"
)

var Messages = driver.Messages

func Dial(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, Option{Messages: Messages}), nil
}
