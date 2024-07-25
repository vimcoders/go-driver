package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go-driver/app/sense/driver"
	"go-driver/app/sense/handler"
	"go-driver/log"
	"go-driver/tcp"

	"gopkg.in/yaml.v3"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "conf", "./sense.conf", "sense.conf")
	flag.Parse()
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	var opt driver.YAML
	if err := yaml.Unmarshal(ymalBytes, &opt); err != nil {
		panic(err.Error())
	}
	handler := handler.MakeHandler(opt)
	addr, err := net.ResolveTCPAddr("tcp4", opt.TCP.Port)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	//tcpAddr := listener.Addr().(*net.TCPAddr)
	ctx, cancel := context.WithCancel(context.Background())
	go tcp.ListenAndServe(ctx, listener, handler)
	log.Infof("running %s", listener.Addr().String())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-quit
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
			log.Info("shutdown ->", s.String())
			cancel()
			listener.Close()
		default:
			log.Info("os.Signal ->", s.String())
			continue
		}
	}
}
