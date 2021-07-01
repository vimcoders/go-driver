package driver

import (
	"database/sql"
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
	Execer
	SetMaxOpenConns(n int)
	Close() error
}

type Scanner interface {
	Scan(scanner func(dest ...interface{}) error) error
}

type Convertor interface {
	Convert() (sql string, args []interface{})
}

type Table interface {
	TableName() string
}

type Marshaler interface {
	Marshal() (str string, err error)
}

type Unmarshaler interface {
	Unmarshal(str string) error
}

type Result sql.Result

type Execer interface {
	Exec(c ...Convertor) error
}
