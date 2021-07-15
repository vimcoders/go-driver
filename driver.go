package driver

import (
	"context"
)

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Close() (err error)
}

type Message interface {
	ToBytes() (b []byte, err error)
}

type Session interface {
	SessionID() int64
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	Send(pkg Message) (err error)
}

type Connector interface {
	Tx(ctx context.Context) (Execer, error)
	SetMaxOpenConns(n int)
	Close() error
}

type Execer interface {
	Exec(i ...interface{}) (interface{}, error)
	ExecContext(ctx context.Context, i ...interface{}) (interface{}, error)
}
