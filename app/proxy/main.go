package main

import (
	"context"
	"flag"
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

	"gopkg.in/yaml.v3"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "conf", "./proxy.conf", "proxy.conf")
	flag.Parse()
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	var opt driver.YAML
	if err := yaml.Unmarshal(ymalBytes, &opt); err != nil {
		panic(err.Error())
	}
	handler := handler.MakeHandler(&opt)
	srv := &http.Server{
		Addr:    opt.HTTP.Port,
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
	s := <-quit
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
		log.Info("os.Signal ->", s.String())
	default:
		log.Info("os.Signal ->", s.String())
		return
	}
	handler.Close()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("Proxy Server Shutdown: %s", err.Error())
	}
}
