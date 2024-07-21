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
	// for i := 1; i < 10000; i++ {
	// 	go func() {
	// 		for {
	// 			handler.LoginRequest()
	// 		}
	// 	}()
	// }
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	cancel()
	listener.Close()
}
