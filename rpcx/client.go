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
	messages  []string
}

func NewClient(c net.Conn) *Client {
	client := &Client{
		w:        c,
		pending:  make(map[uint32]chan Message),
		Buffsize: 16 * 1024,
		Timeout:  time.Second * 120,
	}
	return client
}

func (x *Client) Register(ctx context.Context, messages ...proto.Message) error {
	if len(x.messages) > 0 {
		return errors.New("len(x.ProtoBuf)> 0")
	}
	for i := 0; i < len(messages); i++ {
		x.messages = append(x.messages, string(proto.MessageName(messages[i]).Name()))
	}
	go x.Pull(ctx)
	go x.Keeplive(ctx)
	return nil
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
		response, err := decode(buffer)
		if err != nil {
			return err
		}
		go x.callback(ctx, response)
	}
}

func (x *Client) callback(ctx context.Context, response Message) error {
	seqNumber := response.SeqNumber()
	if seqNumber == 0 {
		var message pb.Message
		if err := proto.Unmarshal(response.Message(), &message); err != nil {
			return err
		}
		x.messages = message.Messages
		return nil
	}
	done := x.done(seqNumber)
	if done == nil {
		return errors.New("done == nil")
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	select {
	case done <- response:
		return nil
	case <-timeoutCtx.Done():
		return nil
	}
}

func (x *Client) Call(ctx context.Context, request proto.Message, reply proto.Message) (err error) {
	for i := 0; i < math.MaxInt32; i++ {
		done, messageId, err := x.Push(request)
		if err != nil {
			log.Error(err.Error())
			continue
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
	return errors.New("try many request")
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
	messageName := string(proto.MessageName(message).Name())
	for i := uint16(0); i < uint16(len(x.messages)); i++ {
		if messageName != x.messages[i] {
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
	for i := uint32(1); i < math.MaxUint16; i++ {
		messageId := x.messageId + i
		if _, ok := x.pending[messageId]; ok {
			continue
		}
		done := make(chan Message, 1)
		x.pending[messageId] = done
		x.messageId = x.messageId%math.MaxUint16 + 1
		return done, messageId
	}
	return nil, 0
}

func (x *Client) done(seqNumber uint32) chan Message {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[seqNumber]; ok {
		delete(x.pending, seqNumber)
		return v
	}
	return nil
}
