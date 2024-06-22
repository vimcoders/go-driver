package log_test

import (
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	syslog "go-driver/log"

	"github.com/golang/glog"
)

func Benchmark_Slog(b *testing.B) {
	logger := slog.New(slog.Default().Handler())
	for i := 0; i < b.N; i++ {
		logger.Info("Go is best language!")
	}
}

func Benchmark_Glog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		glog.Info("Go is best language!")
	}
	glog.Flush()
}

func Benchmark_Syslog(b *testing.B) {
	syslogger := syslog.NewSysLogger()
	for i := 0; i < b.N; i++ {
		syslogger.Error("Go is best language!")
		syslogger.Errorf("Go is best language! %d", time.Now().Unix())
	}
	syslogger.Close()
}

func Benchmark_log(b *testing.B) {
	log.SetPrefix("[MyApp]")
	log.SetOutput(os.Stdout)
	for i := 0; i < b.N; i++ {
		log.Println("Go is best language!")
	}
}
