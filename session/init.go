package session

import (
	"context"
	"net"
	"net/http"

	"github.com/vimcoders/go-driver/log"
	"golang.org/x/net/websocket"
)

var (
	CloseCtx, CloseFunc = context.WithCancel(context.Background())
)

func init_tcp() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	go func() {
		if e := recover(); e != nil {
			log.Error("accept err %v", e)
		}
		for {
			c, err := l.Accept()
			if err != nil {
				log.Error("Accept err %v", err)
				continue
			}
			closeCtx, closeFunc := context.WithCancel(context.Background())
			s := &Session{
				Conn:       c,
				C:          make(chan []byte, 1),
				CancelFunc: closeFunc,
			}
			s.OnMessage = func(b []byte) error {
				return nil
			}
			go s.Pull(closeCtx)
			go s.Push(closeCtx)
		}
	}()
}

func init_websocket() {
	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		closeCtx, closeFunc := context.WithCancel(context.Background())
		s := &Session{
			Conn:       ws,
			C:          make(chan []byte, 1),
			CancelFunc: closeFunc,
		}
		s.OnMessage = func(b []byte) error {
			return nil
		}
		go s.Push(closeCtx)
		s.Pull(closeCtx)
	}))
	go func() {
		if e := recover(); e != nil {
			log.Error("accept rr %v", e)
		}
		if err := http.ListenAndServe(":8889", nil); err != nil {
			log.Error("listent %v", err)
		}
	}()
}
