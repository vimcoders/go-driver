package log_test

import (
	"log"
	"log/slog"
	"os"
	"testing"

	syslog "go-driver/log"

	"github.com/golang/glog"
)

// go test -bench=Benchmark_Slog --benchmem
func Benchmark_Slog(b *testing.B) {
	logger := slog.New(slog.Default().Handler())
	for i := 0; i < b.N; i++ {
		logger.Info("Go is best language!")
	}
}

// go test -bench=Benchmark_Glog --benchmem
func Benchmark_Glog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		glog.Info("Go is best language!")
	}
	glog.Flush()
}

// go test -bench=Benchmark_Syslog --benchmem
func Benchmark_Syslog(b *testing.B) {
	syslogger := syslog.NewSysLogger()
	for i := 0; i < b.N; i++ {
		syslogger.Info("Go is best language!")
	}
	syslogger.Close()
}

// go test -bench=Benchmark_log --benchmem
func Benchmark_log(b *testing.B) {
	log.SetPrefix("[MyApp]")
	log.SetOutput(os.Stdout)
	for i := 0; i < b.N; i++ {
		log.Println("Go is best language!")
	}
}
