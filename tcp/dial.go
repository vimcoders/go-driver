package tcp

import "net"

func Dial(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn), nil
}
