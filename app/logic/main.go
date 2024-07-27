// 逻辑服务 处理请求
package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go-driver/app/logic/driver"
	"go-driver/app/logic/handle"
	"go-driver/log"
)

func main() {
	log.Info("NumCPU: ", runtime.NumCPU())
	ctx, cancel := context.WithCancel(context.Background())
	handler := handle.MakeHandler(driver.ReadOption())
	go handler.ListenAndServe(ctx)
	log.Info("logic running")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-quit
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
			log.Info("shutdown ->", s.String())
			cancel()
			return
		default:
			log.Info("os.Signal ->", s.String())
			continue
		}
	}
}
