package rpcx

import (
	"bufio"
	"context"
	"errors"
	"net"
	"reflect"
	"time"

	"go-driver/driver"
	"go-driver/log"

	"google.golang.org/protobuf/proto"
)

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(ctx context.Context, listener net.Listener, handler driver.Handler) {
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

type Server struct {
	net.Conn
	Handler  interface{} // handler to invoke, http.DefaultServeMux if nil
	Buffsize uint16
	Timeout  time.Duration
}

func (x *Server) Register(handler interface{}) {
	x.Handler = handler
	go x.Poll(context.Background())
}

func (x *Server) ServeRPCX(w ResponsePusher, b []byte, opt Option) (err error) {
	method := reflect.ValueOf(x.Handler).MethodByName(opt.Get(MESSAGENAME))
	// if ok := method.IsNil(); ok {
	// 	return fmt.Errorf("method.IsNil() %s", method)
	// }
	// if ok := method.IsZero(); ok {
	// 	return errors.New("method.IsZero()")
	// }
	t := method.Type()
	if t.NumIn() < 2 {
		return errors.New("t.NumIn() < 2")
	}
	e := t.In(1).Elem()
	in, ok := reflect.New(e).Interface().(proto.Message)
	if !ok {
		return errors.New("!ok")
	}
	if err := proto.Unmarshal(b, in); err != nil {
		return err
	}
	values := method.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(in)})
	if len(values) <= 0 {
		return errors.New("len(values) <= 0")
	}
	response := &Response{
		Option: opt,
		Conn:   x.Conn,
	}
	response.Push(context.Background(), values[0].Interface().(proto.Message))
	return nil
}

func (x *Server) Poll(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
		if err := x.Close(); err != nil {
			log.Error(err.Error())
		}
	}()
	buffer := bufio.NewReaderSize(x.Conn, int(x.Buffsize))
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		message, err := decodeRequest(buffer)
		if err != nil {
			return err
		}
		response := &Response{
			Option: message.Option,
			Conn:   x.Conn,
		}
		if h, ok := x.Handler.(Handler); ok {
			go h.ServeRPCX(response, message.Message, message.Option)
			continue
		}
		go x.ServeRPCX(response, message.Message, message.Option)
	}
}
