package driver

import (
	"bufio"
	"go-driver/session"
)

type Session = session.Session
type Request = session.Request
type Response = session.Response

var Messages = session.Messages

func ReadRequest(b *bufio.Reader) (request Request, err error) {
	return session.ReadRequest(b)
}

func NewRequest(kind uint16) Request {
	return session.NewRequest(kind)
}
