package driver

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

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
	log.Output(2, fmt.Sprintf("\x1b[%dm[Debug] %s\x1b[0m", color_white, fmt.Sprintf(format, v...)))
}

func (log *Syslogger) Info(format string, v ...interface{}) {
	log.Output(2, fmt.Sprintf("\x1b[%dm[Info] %s\x1b[0m", color_green, fmt.Sprintf(format, v...)))
}

func (log *Syslogger) Warning(format string, v ...interface{}) {
	log.Output(2, fmt.Sprintf("\x1b[%dm[Warning] %s\x1b[0m", color_yellow, fmt.Sprintf(format, v...)))
}

func (log *Syslogger) Error(format string, v ...interface{}) {
	log.Output(2, fmt.Sprintf("\x1b[%dm[Error] %s\x0A%s\x1b[0m", color_red, fmt.Sprintf(format, v...), string(debug.Stack())))
}

func (log *Syslogger) Close() error {
	return nil
}

func NewSyslogger() (Logger, error) {
	return &Syslogger{log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)}, nil
}
