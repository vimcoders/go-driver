package rpcx

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"go-driver/log"
	"go-driver/pb"

	"google.golang.org/protobuf/proto"
)

type Client struct {
	w net.Conn
	sync.RWMutex
	messageId uint32
	pending   map[uint32]chan Message
	Buffsize  uint16
	Timeout   time.Duration
	ProtoBuf
}

func NewClient(c net.Conn) *Client {
	client := &Client{
		w:        c,
		pending:  make(map[uint32]chan Message),
		Buffsize: 16 * 1024,
		Timeout:  time.Second * 120,
		ProtoBuf: messages,
	}
	go client.Pull(context.Background())
	go client.Keeplive(context.Background())
	return client
}

func (x *Client) Keeplive(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		if err := x.Ping(ctx); err != nil {
			log.Error(err.Error())
			return
		}
	}
}

func (x *Client) Pull(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
			debug.PrintStack()
		}
		if err := x.w.Close(); err != nil {
			log.Error(err.Error())
		}
	}()
	buffer := bufio.NewReaderSize(x.w, int(x.Buffsize))
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.w.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		message, err := decode(buffer)
		if err != nil {
			return err
		}
		done := x.done(message.PackageNumber())
		if done == nil {
			continue
		}
		done <- message
	}
}

func (x *Client) Call(ctx context.Context, request proto.Message, reply proto.Message) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	done, messageId, err := x.Push(request)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		x.done(messageId)
		return errors.New("timeout")
	case v := <-done:
		close(done)
		return proto.Unmarshal(v.Message(), reply)
	}
}

func (x *Client) Ping(ctx context.Context) (err error) {
	log.Info("ping...")
	var reply pb.PingResponse
	if err := x.Call(ctx, &pb.PingRequest{}, &reply); err != nil {
		return err
	}
	log.Info("ping response ...")
	return nil
}

func (x *Client) Push(message proto.Message) (chan Message, uint32, error) {
	for i := uint16(0); i < uint16(len(x.ProtoBuf)); i++ {
		if proto.MessageName(message) != proto.MessageName(x.ProtoBuf[i]) {
			continue
		}
		ch, messageId := x.newPending()
		if ch == nil {
			return nil, 0, errors.New("ch == nil")
		}
		b, err := encode(messageId, i, message)
		if err != nil {
			return nil, 0, err
		}
		if _, err := x.w.Write(b); err != nil {
			return nil, 0, err
		}
		return ch, messageId, nil
	}
	return nil, 0, fmt.Errorf("message %s not registered", proto.MessageName(message))
}

func (x *Client) newPending() (chan Message, uint32) {
	x.Lock()
	defer x.Unlock()
	for i := uint32(1); i < math.MaxInt32; i++ {
		messageId := x.messageId + i
		if _, ok := x.pending[messageId]; ok {
			continue
		}
		done := make(chan Message, 1)
		x.pending[messageId] = done
		x.messageId = x.messageId%math.MaxInt32 + 1
		return done, messageId
	}
	return nil, 0
}

func (x *Client) done(messageId uint32) chan Message {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[messageId]; ok {
		delete(x.pending, messageId)
		return v
	}
	return nil
}
