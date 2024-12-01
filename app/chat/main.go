package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/vimcoders/go-driver/app/chat/handler"
	"github.com/vimcoders/go-driver/log"
)

func main() {
	log.Info("NumCPU: ", runtime.NumCPU())
	ctx, cancel := context.WithCancel(context.Background())
	handler := handler.MakeHandler(ctx)
	handler.ListenAndServe(ctx)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	s := <-quit
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
		cancel()
		handler.Close()
		log.Info("shutdown ->", s.String())
	default:
		log.Info("os.Signal ->", s.String())
	}
}
