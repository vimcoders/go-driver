// 性能测试
package main

import (
	"go-driver/app/benchmark/grpcx"
	"go-driver/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// func main() {
// 	runtime.GOMAXPROCS(2)
// 	log.Info("runtime.NumCPU: ", runtime.NumCPU())
// 	for i := 0; i < 2; i++ {
// 		client := benchmark.Client{
// 			Url:      "http://127.0.0.1:9800/api/v1/passport/login",
// 			CometUrl: "127.0.0.1:9600",
// 		}
// 		if err := client.Login(); err != nil {
// 			log.Error(err.Error())
// 			continue
// 		}
// 		log.Debug("NumGoroutine", runtime.NumGoroutine())
// 	}
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
// 	ticker := time.NewTicker(time.Second)
// 	for {
// 		select {
// 		case <-quit:
// 			return
// 		case <-ticker.C:
// 			log.Debug("NumGoroutine", runtime.NumGoroutine())
// 		}
// 	}
// }

func main() {
	runtime.GOMAXPROCS(2)
	log.Info("runtime.NumCPU: ", runtime.NumCPU())
	client, err := grpcx.Dial("udp", "127.0.0.1:8972")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 4000; i++ {
		go client.BenchmarkQUIC()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			log.Debug("NumGoroutine", runtime.NumGoroutine())
		}
	}
}
