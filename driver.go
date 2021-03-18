package driver

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Close() (err error)
}

type Message interface {
	ToMessage() (msg []byte, err error)
}

type Session interface {
	SessionID() int64
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	Write(pkg Message) (err error)
}
