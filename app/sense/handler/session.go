package handler

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"time"

	"go-driver/driver"
	"go-driver/grpcx"
	"go-driver/log"

	"google.golang.org/protobuf/proto"
)

type Session struct {
	Id    string
	Token string
	driver.Marshal
	driver.Unmarshal
	iClient *grpcx.Client
	net.Conn
	Buffsize int
	Header   int
	Timeout  time.Duration
}

func (x *Session) Handle(w io.Writer, request []byte) {
	// args, _, err := x.Unmarshal.Unmarshal(request)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return
	// }
	// methodName := proto.MessageName(args).Name()
	// method := reflect.ValueOf(x).MethodByName(string(methodName))
	// values := method.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(args)})
	// if len(values) <= 0 {
	// 	log.Error("len(values) <= 0")
	// 	return
	// }
	// x.Push(context.Background(), values[0].Interface().(proto.Message))
}

func (x *Session) Poll(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
		}
		if err != nil {
			log.Error(err.Error())
			debug.PrintStack()
		}
		x.Close()
	}()
	buffer := bufio.NewReaderSize(x.Conn, x.Buffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		headerBytes, err := buffer.Peek(x.Header)
		if err != nil {
			return err
		}
		header := binary.BigEndian.Uint32(headerBytes)
		log.Debug(headerBytes, x.Header, header)
		if int(header) > buffer.Size() {
			return fmt.Errorf("header %v too long", header)
		}
		message, err := buffer.Peek(int(header) + len(headerBytes))
		if err != nil {
			return err
		}
		if len(message) < x.Header {
			return errors.New("len(bodyBytes) < Header+Proto")
		}
		x.Handle(x.Conn, message[x.Header:])
		if _, err := buffer.Discard(len(message)); err != nil {
			return err
		}
	}
}

func (x *Session) Push(ctx context.Context, message proto.Message) error {
	b, err := x.Marshal.Marshal(message)
	if err != nil {
		return err
	}
	var buffer = driver.NewBuffer(4)
	binary.BigEndian.PutUint32(buffer[:], uint32(len(b)))
	buffer.Write(b)
	if _, err := x.Write(buffer); err != nil {
		return err
	}
	return nil
}

func (x *Session) Close() error {
	return nil
}
