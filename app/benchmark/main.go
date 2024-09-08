// 性能测试
package main

import (
	"context"
	"crypto/tls"
	benchmark "go-driver/app/benchmark/client"
	"go-driver/driver"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/quicx"
	"go-driver/tcp"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	_ "net/http/pprof"
)

// func main() {
// 	runtime.GOMAXPROCS(2)
// 	log.Info("runtime.NumCPU: ", runtime.NumCPU())
// 	for i := 0; i < 2; i++ {
// 		client := benchmark.Client{
// 			Url:      "http://127.0.0.1:9800/api/v1/passport/login",
// 			CometUrl: "127.0.0.1:9600",
// 		}
// 		if err := client.Login(); err != nil {
// 			log.Error(err.Error())
// 			continue
// 		}
// 		log.Debug("NumGoroutine", runtime.NumGoroutine())
// 	}
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
// 	ticker := time.NewTicker(time.Second)
// 	for {
// 		select {
// 		case <-quit:
// 			return
// 		case <-ticker.C:
// 			log.Debug("NumGoroutine", runtime.NumGoroutine())
// 		}
// 	}
// }

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	log.Info("runtime.NumCPU: ", runtime.NumCPU())
	conn, err := quicx.Dial("127.0.0.1:9700", &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
		MaxVersion:         tls.VersionTLS13,
	}, &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	bot := benchmark.Client{
		ServiceDesc: pb.Parkour_ServiceDesc,
		Client:      tcp.NewClient(conn, tcp.Option{}),
		Pool: sync.Pool{
			New: func() any {
				return &driver.Message{}
			},
		},
	}
	bot.Client.Register(&bot)
	for i := 0; i < 20000; i++ {
		go bot.Ping(context.Background())
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			log.Debug(" NumGoroutine ", runtime.NumGoroutine())
		}
	}
}
