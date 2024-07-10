package tcp

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	"go-driver/driver"
	"go-driver/log"
)

// Config stores tcp server properties
type Config struct {
	Address    string        `yaml:"address"`
	MaxConnect uint32        `yaml:"max-connect"`
	Timeout    time.Duration `yaml:"timeout"`
	Key        string        `yaml:"key"`
}

// ListenAndServeWithSignal binds port and handle requests, blocking until receive stop signal
func ListenAndServeWithSignal(cfg *Config, handler driver.Handler) error {
	// closeChan := make(chan struct{})
	// sigCh := make(chan os.Signal)
	// signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	// go func() {
	// 	sig := <-sigCh
	// 	switch sig {
	// 	case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
	// 		closeChan <- struct{}{}
	// 	}
	// }()
	// listener, err := net.Listen("tcp", cfg.Address)
	// if err != nil {
	// 	return err
	// }
	// //cfg.Address = listener.Addr().String()
	// //logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	// ListenAndServe(listener, handler, closeChan)
	return nil
}

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(ctx context.Context, listener net.Listener, handler driver.Handler) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		log.Debugf("new conn %s", conn.RemoteAddr().String())
		handler.Handle(ctx, conn)
	}
}
