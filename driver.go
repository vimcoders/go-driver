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
	Write(pkg Message) (err error)
}

type Connector interface {
	Tx(ctx context.Context) (Tx, error)
	Conn(ctx context.Context) (Conn, error)
	SetMaxOpenConns(n int)
	Close() error
}

type Tx interface {
	Exec(obj Object) (interface{}, error)
	ExecContext(ctx context.Context, obj Object) (interface{}, error)
	Rollback() (err error)
	Commit() (err error)
}

type Conn interface {
	Exec(obj Object) (interface{}, error)
	ExecContext(ctx context.Context, obj Object) (interface{}, error)
	Tx(ctx context.Context) (Tx, error)
	Close() (err error)
}

type Object interface {
	Table
	Convert() (sql string, args []interface{})
	Scan(scanner func(dest ...interface{}) error) error
}

type Table interface {
	TableName() string
}

type Marshaler interface {
	Marshal(str string, err error)
}

type Unmarshaler interface {
	Unmarshal(str string) error
}
