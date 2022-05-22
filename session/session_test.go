package session

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	m.Run()
	fmt.Println("end")
}

// 测试发送消息
func TestHelloWorld(t *testing.T) {
	init_tcp()
	var waitGroup sync.WaitGroup
	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		c, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			t.Error(err)
			return
		}
		s := &Session{
			Conn: c,
			C:    make(chan []byte, 1),
		}
		s.OnMessage = func(b []byte) error {
			t.Log(string(b))
			return nil
		}
		go func() {
			if err := s.Pull(); err != nil {
				t.Log(err)
			}
		}()
		go func() {
			if err := s.Push(); err != nil {
				t.Log(err)
			}
		}()
		go func() {
			defer waitGroup.Done()
			for k := 0; k < 100; k++ {
				s.C <- []byte("hello")
				time.Sleep(time.Second)
			}
		}()
	}
	waitGroup.Wait()
}

// 测试发送消息
func TestWebSocket(t *testing.T) {
	init_websocket()
	var waitGroup sync.WaitGroup
	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		ws, err := websocket.Dial("ws://localhost:8889/ws", "", "http://localhost/")
		if err != nil {
			t.Error(err)
			return
		}
		s := &Session{
			Conn: ws,
			C:    make(chan []byte, 1),
		}
		s.OnMessage = func(b []byte) error {
			t.Log(string(b))
			return nil
		}
		go s.Pull()
		go s.Push()
		go func() {
			defer waitGroup.Done()
			defer s.Close()
			for k := 0; k < 100; k++ {
				s.C <- []byte("hello")
				time.Sleep(time.Second)
			}
		}()
	}
	waitGroup.Wait()
}
