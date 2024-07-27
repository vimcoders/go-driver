package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"go-driver/app/proxy/driver"
	"go-driver/app/proxy/handler"
	"go-driver/app/proxy/router"
	"go-driver/log"
)

func main() {
	option := driver.ReadOption()
	handler := handler.MakeHandler(option)
	srv := &http.Server{
		Addr:    option.HTTP.Port,
		Handler: router.NewRouter(handler),
	}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Error(fmt.Sprintf("%s", e))
				debug.PrintStack()
			}
		}()
		if err := srv.ListenAndServe(); err != nil {
			log.Errorf("listen: %s", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-quit
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
			log.Info("shutdown ->", s.String())
			handler.Close()
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Errorf("Proxy Server Shutdown: %s", err.Error())
			}
		default:
			log.Info("os.Signal ->", s.String())
			continue
		}
	}
}
