package main

import (
	"go-driver/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	log.Info("NumCPU: ", runtime.NumCPU())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	s := <-quit
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
		log.Info("shutdown ->", s.String())
	default:
		log.Info("os.Signal ->", s.String())
	}
}
