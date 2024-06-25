package rpcx

import (
	"bufio"
	"context"
	"errors"
	"go-driver/log"
	"go-driver/pb"
	"net"
	"reflect"
	"time"

	"google.golang.org/protobuf/proto"
)

type Handle struct {
	net.Conn
	Handler  interface{} // handler to invoke, http.DefaultServeMux if nil
	Buffsize uint16
	Timeout  time.Duration
}

func (x *Handle) Register(handler interface{}) {
	x.Handler = handler
	go x.Poll(context.Background())
}

func (x *Handle) ServeRPCX(w ResponsePusher, b []byte, opt Option) (err error) {
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
	pusher := &Pusher{
		Option: opt,
		Conn:   x.Conn,
	}
	pusher.Push(context.Background(), values[0].Interface().(proto.Message))
	return nil
}

func (x *Handle) Poll(ctx context.Context) (err error) {
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
		message, err := decode(buffer)
		if err != nil {
			return err
		}
		pusher := &Pusher{
			Option: message.Option,
			Conn:   x.Conn,
		}
		if h, ok := x.Handler.(Handler); ok {
			go h.ServeRPCX(pusher, message.Message, message.Option)
			continue
		}
		go x.ServeRPCX(pusher, message.Message, message.Option)
	}
}

func (x *Handle) Push(ctx context.Context, iMessage *pb.Message) error {
	pusher := Pusher{
		Conn:    x.Conn,
		Timeout: time.Second * 120,
	}
	return pusher.Push(ctx, iMessage)
}

type Pusher struct {
	Option
	net.Conn
	Timeout time.Duration
}

func (x *Pusher) Push(ctx context.Context, iMessage proto.Message) error {
	b, err := proto.Marshal(iMessage)
	if err != nil {
		return err
	}
	message := &pb.Message{
		Message: b,
	}
	message.Option = append(message.Option, &pb.Option{Key: MESSAGEID, Value: x.Get(MESSAGEID)})
	response, err := encode(message)
	if err != nil {
		return err
	}
	if err := x.SetWriteDeadline(time.Now().Add(x.Timeout)); err != nil {
		return err
	}
	if _, err := x.Conn.Write(response); err != nil {
		return err
	}
	return nil
}
