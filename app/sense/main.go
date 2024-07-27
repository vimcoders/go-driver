package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go-driver/app/sense/driver"
	"go-driver/app/sense/handler"
	"go-driver/log"
)

func main() {
	log.Info("NumCPU: ", runtime.NumCPU())
	quit := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	handler := handler.MakeHandler(driver.ReadOption())
	go handler.ListenAndServe(ctx)
	log.Info("sense running")
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
