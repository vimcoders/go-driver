package log

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

var logger = NewSysLogger()

func Debug(a ...any) {
	logger.Debug(a...)
}

func Debugf(format string, a ...any) {
	logger.Debugf(format, a...)
}

func Info(a ...any) {
	logger.Info(a...)
}

func Infof(format string, a ...any) {
	logger.Infof(format, a...)
}

func Warn(a ...any) {
	logger.Warn(a...)
}

func Warnf(format string, a ...any) {
	logger.Warnf(format, a...)
}

func Error(a ...any) {
	logger.Error(a...)
}

func Errorf(format string, a ...any) {
	logger.Errorf(format, a...)
}

type SysLogger struct {
	Handler
}

func NewSysLogger() *SysLogger {
	return &SysLogger{
		Handler: NewSyncBuffer(),
	}
}

var pool sync.Pool = sync.Pool{
	New: func() any {
		return &Buffer{}
	},
}

func (x *SysLogger) Debug(a ...any) {
	x.log(" DEBUG ", a...)
}

func (x *SysLogger) Debugf(format string, a ...any) {
	x.Debug(fmt.Sprintf(format, a...))
}

func (x *SysLogger) Info(a ...any) {
	x.log(" INFO ", a...)
}

func (x *SysLogger) Infof(format string, a ...any) {
	x.Info(fmt.Sprintf(format, a...))
}

func (x *SysLogger) Error(a ...any) {
	x.log(" ERROR ", a...)
}

func (x *SysLogger) Errorf(format string, a ...any) {
	x.Error(fmt.Sprintf(format, a...))
}

func (x *SysLogger) Warn(a ...any) {
	x.log(" WARN ", a...)
}

func (x *SysLogger) Warnf(format string, a ...any) {
	x.Warn(fmt.Sprintf(format, a...))
}

func (x *SysLogger) log(prefix string, a ...any) {
	buffer := pool.Get().(*Buffer)
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05 "))
	buffer.WriteString(prefix)
	x.Handler.Handle(context.Background(), fmt.Appendln(*buffer, a...))
	buffer.Reset()
	pool.Put(buffer)
}

func (x *SysLogger) Close() error {
	return x.Handler.Close()
}

type SyncHandler struct {
}

func NewSyncBuffer() Handler {
	return &SyncHandler{}
}

func (x *SyncHandler) Handle(ctx context.Context, b []byte) error {
	if _, err := os.Stdout.Write(b); err != nil {
		return err
	}
	return nil
}

func (x *SyncHandler) Close() error {
	return nil
}
