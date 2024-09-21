package handler

import (
	"context"
	"encoding/json"
	"flag"
	"go-driver/app/chat/driver"
	"go-driver/etcdx"
	"go-driver/grpcx"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/quicx"
	"net"
	"os"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

type Option = driver.Option

type Handler struct {
	Option
	*etcd.Client
	pb.ChatServer
}

func MakeHandler(ctx context.Context) *Handler {
	h := &Handler{}
	var fileName string
	flag.StringVar(&fileName, "option", "chat.conf", "chat.conf")
	flag.Parse()
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(ymalBytes, &h.Option); err != nil {
		panic(err)
	}
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{h.Etcd.Endpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	h.Client = cli
	return h
}

func (x *Handler) ListenAndServe(ctx context.Context) {
	// addr, err := net.ResolveTCPAddr("tcp4", opt.Addr.Port)
	// if err != nil {
	// 	panic(err)
	// }
	// listener, err := net.ListenTCP("tcp", addr)
	listener, err := quicx.Listen("udp", x.QUIC.LocalAddr, GenerateTLSConfig(), &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	go grpcx.ListenAndServe(ctx, listener, x)
	b, err := json.Marshal(&etcdx.Service{
		Kind:      "Chat",
		Internet:  x.QUIC.Internet,
		LocalAddr: x.QUIC.LocalAddr,
		Network:   "QUIC",
	})
	if err != nil {
		panic(err)
	}
	if _, err := x.Client.Put(ctx, x.Etcd.Join("chat"), string(b)); err != nil {
		panic(err)
	}
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, conn net.Conn) {
	log.Infof("new conn %s", conn.RemoteAddr().String())
	cli := grpcx.NewClient(conn, grpcx.Option{ServiceDesc: pb.Chat_ServiceDesc})
	if err := cli.Register(ctx, x); err != nil {
		log.Error(err.Error())
	}
}
