package log

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

func buildf(depth int, prefix string, format string, a ...any) Buffer {
	buffer := pool.Get().(*Buffer)
	_, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		file = "???"
	}
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05 "))
	buffer.WriteString(filepath.Base(file))
	buffer.WriteString(fmt.Sprintf(":%v", line))
	buffer.WriteString(prefix)
	fmt.Fprintf(buffer, format, a...)
	buffer.WriteString("\n")
	return *buffer
}

func build(depth int, prefix string, a ...any) Buffer {
	var buffer Buffer
	_, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		file = "???"
	}
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05 "))
	buffer.WriteString(filepath.Base(file))
	buffer.WriteString(fmt.Sprintf(":%v", line))
	buffer.WriteString(prefix)
	fmt.Fprint(&buffer, a...)
	buffer.WriteString("\n")
	return buffer
}

func (x *SysLogger) Debug(a ...any) {
	buffer := build(1+1, " DEBUG ", a...)
	x.Handler.Handle(context.Background(), buffer)
}

func (x *SysLogger) Debugf(format string, a ...any) {
	buffer := buildf(1+1, " DEBUG ", format, a...)
	x.Handler.Handle(context.Background(), buffer)
}

func (x *SysLogger) Info(a ...any) {
	buffer := build(1+1, " INFO ", a...)
	x.Handler.Handle(context.Background(), buffer)
}

func (x *SysLogger) Infof(format string, a ...any) {
	buffer := buildf(1+1, " INFO ", format, a...)
	defer pool.Put(&buffer)
	x.Handler.Handle(context.Background(), buffer)
}

func (x *SysLogger) Error(a ...any) {
	buffer := build(1+1, " ERROR ", a...)
	x.Handler.Handle(context.Background(), buffer)
}

func (x *SysLogger) Errorf(format string, a ...any) {
	buffer := buildf(1+1, " ERROR ", format, a...)
	x.Handler.Handle(context.Background(), buffer)
}

func (x *SysLogger) Warn(a ...any) {
	buffer := build(1+1, " WARN ", a...)
	x.Handler.Handle(context.Background(), buffer)
}

func (x *SysLogger) Warnf(format string, a ...any) {
	buffer := buildf(1+1, " WARN ", format, a...)
	x.Handler.Handle(context.Background(), buffer)
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
