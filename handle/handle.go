package handle

import (
	"bufio"
	"context"
	"errors"
	"go-driver/driver"
	"go-driver/log"
	"net"
	"runtime/debug"
	"time"

	"google.golang.org/protobuf/proto"
)

// 一个tcp或者udp的解析器，它的主要工作是解析操作系统从网卡上获取到的二进制数据流
type Handle struct {
	w        net.Conn
	Buffsize int
	Timeout  time.Duration
	Handler
	driver.Marshal
	driver.Unmarshal
}

// 从一个tcp或者udp连接构造一个解析器
func NewHandle(w net.Conn) *Handle {
	return &Handle{
		w:         w,
		Buffsize:  512,
		Timeout:   time.Minute,
		Marshal:   Messages,
		Unmarshal: Messages,
	}
}

// 这个解析器将从这里开始工作
func (x *Handle) Pull(ctx context.Context) (err error) {
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
	buffer := bufio.NewReaderSize(x.w, x.Buffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.w.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		req, err := decode(buffer)
		if err != nil {
			return err
		}
		if x.Handler == nil {
			continue
		}
		request, err := x.Unmarshal.Unmarshal(req)
		if err != nil {
			return err
		}
		reply, err := x.Unmarshal.Unmarshal(NewRequest(req.Reply()))
		if err != nil {
			return err
		}
		if err := x.Handler.Handle(ctx, request, reply); err != nil {
			return err
		}
		if err := x.Push(ctx, reply); err != nil {
			return err
		}
	}
}

// 我们将会向网卡发送一段二进制流，告诉对方我们处理二进制的结果
func (x *Handle) Push(ctx context.Context, message proto.Message) error {
	response, err := x.Marshal.Marshal(message)
	if err != nil {
		return err
	}
	if _, err := x.w.Write(response); err != nil {
		return err
	}
	return nil
}

// 我们将在这里关闭一个tcp或者udp连接
func (x *Handle) Close() error {
	return x.w.Close()
}
