package log_test

import (
	"log/slog"
	"testing"

	"github.com/vimcoders/go-driver/log"
)

func Benchmark_Slog(b *testing.B) {
	for i := 0; i < 10000; i++ {
		slog.Info("Go is best language!")
	}
}

func Benchmark_Liblog(b *testing.B) {
	for i := 0; i < 100000; i++ {
		log.Info("Go is best language!")
	}
}

func Benchmark_Glog(b *testing.B) {
	// for i := 0; i < 100000; i++ {
	// 	glog.Info("Go is best language!")
	// }
	// glog.Flush()
}

func Benchmark_Syslog(b *testing.B) {
	syslogger := log.NewSysLogger()
	for i := 0; i < 10; i++ {
		syslogger.Error("Go is best language!")
	}
	syslogger.Close()
}
