// 性能测试
package main

import (
	benchmark "go-driver/app/benchmark/client"
	"go-driver/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(2)
	for i := 0; i < 1000; i++ {
		client := benchmark.Client{
			Url:      "http://127.0.0.1:9800/api/v1/passport/login",
			CometUrl: "127.0.0.1:9600",
		}
		if err := client.Login(); err != nil {
			log.Error(err.Error())
			continue
		}
		log.Debug("NumGoroutine", runtime.NumGoroutine())
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
