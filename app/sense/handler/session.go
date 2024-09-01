package handler

import (
	"context"
	"errors"
	"go-driver/grpcx"
	"go-driver/tcp"
)

type Session struct {
	UserId int64
	Level  int16
	Name   string
	Icon   string
	tcp.Client
	iClient grpcx.Client
	s       Scene
}

func (x *Session) ServeTCP(ctx context.Context, request []byte) error {
	return nil
}

func (x *Session) Bind(s Scene) error {
	if x.s != nil {
		return errors.New("x.s != nil")
	}
	x.s = s
	return nil
}

func (x *Session) Close() error {
	return nil
}
