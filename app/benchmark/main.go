// 性能测试
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/pb"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
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

type Handle struct {
	total uint64
	unix  int64
	pb.ParkourServer
}

// MakeHandler creates a Handler instance
func MakeHandler() *Handle {
	return &Handle{}
}

// Handle receives and executes redis commands
func (x *Handle) Handle(ctx context.Context, conn net.Conn) {
	cli := grpcx.NewClient(conn, grpcx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	//go cli.Keeplive(ctx, &pb.PingRequest{})
	cli.Register(x)
}

func (x *Handle) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	unix := time.Now().Unix()
	x.total++
	if unix != x.unix {
		fmt.Println(x.total, " request/s", " NumGoroutine ", runtime.NumGoroutine())
		x.total = 0
		x.unix = unix
	}
	return &pb.PingResponse{Message: req.Message}, nil
}

func (x *Handle) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{Code: http.StatusOK}, nil
}

// Close stops handler
func (x *Handle) Close() error {
	return nil
}

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	log.Info("runtime.NumCPU: ", runtime.NumCPU())
	conn, err := tls.Dial("tcp", "127.0.0.1:9600", &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
		MaxVersion:         tls.VersionTLS13,
	})
	if err != nil {
		panic(err)
	}
	clientInterface := grpcx.NewClient(conn, grpcx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	clientInterface.Register(MakeHandler())
	client := pb.NewParkourClient(clientInterface)
	for i := 0; i < 100; i++ {
		go func() {
			for {
				client.Ping(context.Background(), &pb.PingRequest{})
			}
		}()
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
