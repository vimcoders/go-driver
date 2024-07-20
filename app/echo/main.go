package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go http.ListenAndServe(":8080", nil)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
