package driver

import (
	"go-driver/session"
	"net"
)

type Session = session.Session
type Request = session.Request
type Response = session.Response

var Messages = session.Messages

func NewRequest(kind uint16) Request {
	return session.NewRequest(kind)
}

func NewSession(w net.Conn) *session.Session {
	return session.NewSession(w)
}
