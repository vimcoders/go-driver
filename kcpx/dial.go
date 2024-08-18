package kcpx

import (
	"go-driver/driver"

	"github.com/xtaci/kcp-go"
)

var Messages = driver.Messages

func Dial(addr string) (Client, error) {
	kcpconn, err := kcp.DialWithOptions(addr, nil, 10, 3)
	if err != nil {
		return nil, err
	}
	return NewClient(kcpconn, Option{Messages: Messages}), nil
}
