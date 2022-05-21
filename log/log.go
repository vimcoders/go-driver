package log

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/vimcoders/go-driver"
)

var logger = NewSyslogger()

func Error(format string, v ...interface{}) {
	logger.Error(format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

func Warning(format string, v ...interface{}) {
	logger.Warning(format, v...)
}

const (
	color_black = uint8(iota + 90)
	color_red
	color_green
	color_yellow
	color_blue
	color_magenta
	color_cyan
	color_white
)

type Syslogger struct {
	*log.Logger
}

func (log *Syslogger) Debug(format string, v ...interface{}) {
	log.Output(3, fmt.Sprintf("\x1b[%dm[debug] %s\x1b[0m", color_white, fmt.Sprintf(format, v...)))
}

func (log *Syslogger) Info(format string, v ...interface{}) {
	log.Output(3, fmt.Sprintf("\x1b[%dm[info] %s\x1b[0m", color_green, fmt.Sprintf(format, v...)))
}

func (log *Syslogger) Warning(format string, v ...interface{}) {
	log.Output(3, fmt.Sprintf("\x1b[%dm[warning] %s\x1b[0m", color_yellow, fmt.Sprintf(format, v...)))
}

func (log *Syslogger) Error(format string, v ...interface{}) {
	log.Output(3, fmt.Sprintf("\x1b[%dm[error] %s\x0A%s\x1b[0m", color_red, fmt.Sprintf(format, v...), string(debug.Stack())))
}

func (log *Syslogger) Close() error {
	return nil
}

func NewSyslogger() driver.Logger {
	return &Syslogger{log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)}
}
