package kcpx

import (
	"github.com/xtaci/kcp-go"
)

func Dial(addr string) (Client, error) {
	kcpconn, err := kcp.DialWithOptions(addr, nil, 10, 3)
	if err != nil {
		return nil, err
	}
	return NewClient(kcpconn, Option{}), nil
}
