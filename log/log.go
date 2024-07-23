package log

import (
	"context"
	"os"
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

func (x *SysLogger) Debug(a ...any) {
	x.log(" [DEBUG] ", a...)
}

func (x *SysLogger) Debugf(format string, a ...any) {
	x.logf(" [DEBUG] ", format, a...)
}

func (x *SysLogger) Info(a ...any) {
	x.log(" [INFO] ", a...)
}

func (x *SysLogger) Infof(format string, a ...any) {
	x.logf(" [INFO] ", format, a...)
}

func (x *SysLogger) Error(a ...any) {
	x.log(" [ERROR] ", a...)
}

func (x *SysLogger) Errorf(format string, a ...any) {
	x.logf(" [ERROR] ", format, a...)
}

func (x *SysLogger) Warn(a ...any) {
	x.log(" [WARN] ", a...)
}

func (x *SysLogger) Warnf(format string, a ...any) {
	x.logf(" [WARN] ", format, a...)
}

func (x *SysLogger) log(prefix string, a ...any) {
	buffer := newPrinter(prefix, a...)
	x.Handler.Handle(context.Background(), *buffer)
	buffer.free()
}

func (x *SysLogger) logf(prefix, format string, a ...any) {
	buffer := newPrinterf(prefix, format, a...)
	x.Handler.Handle(context.Background(), *buffer)
	buffer.free()
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
