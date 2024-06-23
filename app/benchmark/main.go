// 性能测试
package main

import (
	"fmt"
	benchmark "go-driver/app/benchmark/client"
	"go-driver/log"
	"go-driver/session"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	var count int64
	for i := 0; i < 8000*4; i++ {
		client := benchmark.Client{
			Url:       "http://127.0.0.1:9800/api/v1/passport/login",
			CometUrl:  "127.0.0.1:9600",
			Marshal:   session.Messages,
			Unmarshal: session.Messages,
		}
		if err := client.Login(); err != nil {
			log.Error(err.Error())
			continue
		}
		count++
		time.Sleep(time.Millisecond * 10)
	}
	log.Info(count)
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
