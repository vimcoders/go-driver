package main

import (
	"context"
	"go-driver/app/chat/handler"
	"go-driver/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	log.Info("NumCPU: ", runtime.NumCPU())
	ctx, cancel := context.WithCancel(context.Background())
	handler := handler.MakeHandler(ctx)
	go handler.ListenAndServe(ctx)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	s := <-quit
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
		cancel()
		log.Info("shutdown ->", s.String())
	default:
		log.Info("os.Signal ->", s.String())
	}
}
