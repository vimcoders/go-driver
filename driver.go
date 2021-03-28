package driver

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Close() (err error)
}

type Message interface {
	Header() []byte
	Length() int32
	Payload() []byte
}

type Session interface {
	SessionID() int64
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	Write(pkg Message) (err error)
}

type Reader interface {
	Read() (Message, error)
}

type Writer interface {
	Write(msg Message) error
}
