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
	Version() uint8
	Protocol() uint16
	Payload() []byte
}

type Session interface {
	SessionID() string
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	Send(msg Message)
}
