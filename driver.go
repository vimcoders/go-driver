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
	Update(i interface{}) (interface{}, error)
}

type Deleter interface {
	Delete(i interface{}) (interface{}, error)
}

type Inserter interface {
	Insert(i interface{}) (interface{}, error)
}

type Queryer interface {
	Query(i interface{}) ([]interface{}, error)
}

type Tx interface {
	Updater
	Deleter
	Inserter
	Queryer
	Commit() error
	Rollback() error
}
