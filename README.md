# go-driver
日志
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Close() (err error)
}
网络包
type Message interface {
	ToBytes() (b []byte, err error)
}
面向连接
type Session interface {
	SessionID() int64
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	Writer
}

type Writer interface {
	Write(pkg Message) (err error)
}
