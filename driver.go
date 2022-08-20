package driver

import (
	"context"
	"io"
)

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Close() (err error)
}

type Connector interface {
	Tx(ctx context.Context) (Tx, error)
	Execer(ctx context.Context) (Execer, error)
	SetMaxOpenConns(n int)
	io.Closer
}

type Table interface {
	TableName() string
}

type Updater interface {
	Update(ctx context.Context, i interface{}) (interface{}, error)
}

type Deleter interface {
	Delete(ctx context.Context, i interface{}) (interface{}, error)
}

type Inserter interface {
	Insert(ctx context.Context, i interface{}) (interface{}, error)
}

type Queryer interface {
	Query(ctx context.Context, i interface{}) ([]interface{}, error)
}

type Tx interface {
	Execer
	Close(ctx context.Context) error
}

type Execer interface {
	Updater
	Deleter
	Inserter
	Queryer
	Close(ctx context.Context) error
}
