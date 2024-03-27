package log

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/vimcoders/go-driver/driver"
)

var logger *SysLogger = NewSysLogger()

func Debug(v ...any) {
	logger.Debug(v...)
}

func Debugf(format string, args ...any) {
	logger.Debug(fmt.Sprintf(format, args...))
}

func Info(v ...any) {
	logger.Info(v...)
}

func Infof(format string, args ...any) {
	logger.Info(fmt.Sprintf(format, args...))
}

func Warn(v ...any) {
	logger.Warn(v...)
}

func Warnf(format string, args ...any) {
	logger.Warn(fmt.Sprintf(format, args...))
}

func Error(v ...any) {
	logger.Error(v...)
}

func Errorf(format string, a ...any) {
	logger.Error(fmt.Sprintf(format, a...))
}

func Connect(token, secret string) {
	//logger.Connect(token, secret)
}

type Handler interface {
	Handle(context.Context, []byte) error
	Close() error
}

type SysLogger struct {
	Handler
	driver.Buffer
	sync.RWMutex
}

func NewSysLogger() *SysLogger {
	return &SysLogger{
		Handler: NewSyncBuffer(),
		Buffer:  make(driver.Buffer, 0, 128),
	}
}

func (x *SysLogger) Output(b ...string) {
	now := time.Now()
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
	}
	x.Lock()
	defer x.Unlock()
	x.Buffer.Reset()
	x.Buffer.WriteString(now.Format("2006-01-02 15:04:05 "))
	x.Buffer.WriteString(file)
	x.Buffer.WriteString(fmt.Sprintf(":%v", line))
	for i := 0; i < len(b); i++ {
		x.Buffer.WriteString(b[i])
	}
	x.Buffer.WriteString("\n")
	x.Handle(context.Background(), x.Buffer)
}

func (x *SysLogger) Debug(v ...any) {
	x.Output(" DEBUG ", fmt.Sprint(v...))
}

func (x *SysLogger) Info(v ...any) {
	x.Output(" INFO ", fmt.Sprint(v...))
}

func (x *SysLogger) Error(v ...any) {
	x.Output(" ERROR ", fmt.Sprint(v...))
}

func (x *SysLogger) Warn(v ...any) {
	x.Output(" WARN ", fmt.Sprint(v...))
}

func (x *SysLogger) Close() error {
	x.Lock()
	defer x.Unlock()
	return x.Handler.Close()
}

type SyncBuffer struct {
}

func NewSyncBuffer() Handler {
	return &SyncBuffer{}
}

func (x *SyncBuffer) Handle(ctx context.Context, b []byte) error {
	if _, err := os.Stdout.Write(b); err != nil {
		return err
	}
	return nil
}

func (x *SyncBuffer) Close() error {
	return nil
}
