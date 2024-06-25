package driver

import (
	"go-driver/handle"
	"net"
)

type Handle = handle.Handle
type Request = handle.Request

var Messages = handle.Messages

func NewRequest(kind uint16) Request {
	return handle.NewRequest(kind)
}

func NewHandle(w net.Conn) *handle.Handle {
	return handle.NewHandle(w)
}
