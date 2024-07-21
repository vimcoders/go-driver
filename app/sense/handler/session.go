package handler

import (
	"io"
)

type Session struct {
}

func (x *Session) Handle(w io.Writer, request []byte) {
}

func (x *Session) Close() error {
	return nil
}
