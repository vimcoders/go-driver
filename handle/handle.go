package handle

import (
	"bufio"
	"context"
	"errors"
	"go-driver/log"
	"net"
	"runtime/debug"
	"time"
)

// 一个tcp或者udp的解析器，它的主要工作是解析操作系统从网卡上获取到的二进制
type Handle struct {
	w        net.Conn
	Buffsize int
	Timeout  time.Duration
	Handler
}

// 从一个tcp或者udp连接构造一个解析器
func NewHandle(w net.Conn) *Handle {
	return &Handle{w: w, Buffsize: 512, Timeout: time.Minute}
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
		// 这里我们将会解析二进制流
		request, err := decode(buffer)
		if err != nil {
			return err
		}
		// 调用接口来处理二进制流
		if err := x.Handle(ctx, request); err != nil {
			return err
		}
	}
}

// 我们将会向网卡发送一段二进制流，告诉对方我们处理二进制的结果
func (x *Handle) Push(ctx context.Context, response []byte) (int, error) {
	return x.w.Write(response)
}

// 我们将在这里关闭一个tcp或者udp连接
func (x *Handle) Close() error {
	return nil
}
