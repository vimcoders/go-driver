// 性能测试
package main

import (
	"fmt"
	benchmark "go-driver/app/benchmark/client"
	"go-driver/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	log.Info(time.Now())
	var count int64
	for i := 0; i < 20000; i++ {
		client := benchmark.Client{
			Url:      "http://127.0.0.1:9800/api/v1/passport/login",
			CometUrl: "127.0.0.1:9600",
		}
		if err := client.Login(); err != nil {
			log.Error(err.Error())
			continue
		}
		count++
		time.Sleep(time.Millisecond * 10)
	}
	log.Info(count, time.Now())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			fmt.Println(runtime.NumGoroutine())
		}
	}
}
