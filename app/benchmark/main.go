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

	_ "net/http/pprof"
)

func main() {
	log.Info("runtime.NumCPU: ", runtime.NumCPU())
	client, err := grpcx.Dial("tcp", "127.0.0.1:9600")
	//client, err := tcpx.Dial("tcp", "127.0.0.1:9600")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
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
			log.Debug(" NumGoroutine ", runtime.NumGoroutine())
		}
	}
}
