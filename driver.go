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

type Connector interface {
	Tx(ctx context.Context) (Tx, error)
	SetMaxOpenConns(n int)
	Close() error
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
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Execer interface {
	Updater
	Deleter
	Inserter
	Queryer
}
