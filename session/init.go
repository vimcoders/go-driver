package session

import (
	"context"
	"net"
	"net/http"

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
		for {
			c, err := l.Accept()
			if err != nil {
				continue
			}
			s := &Session{
				Conn: c,
				C:    make(chan []byte, 1),
			}
			s.OnMessage = func(b []byte) error {
				s.C <- []byte("response")
				return nil
			}
			go func() {
				if err := s.Pull(); err != nil {
					panic(err)
				}
			}()
			go func() {
				if err := s.Push(); err != nil {
					panic(err)
				}
			}()
		}
	}()
}

func init_websocket() {
	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		s := &Session{
			Conn: ws,
			C:    make(chan []byte, 1),
		}
		s.OnMessage = func(b []byte) error {
			s.C <- []byte("response")
			return nil
		}
		go func() {
			if err := s.Push(); err != nil {
				panic(err)
			}
		}()
		if err := s.Pull(); err != nil {
			panic(err)
		}
	}))
	go func() {
		http.ListenAndServe(":8889", nil)
	}()
}
