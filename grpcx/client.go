package grpcx

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"path/filepath"
	"sync"
	"time"

	"go-driver/log"
	"go-driver/pb"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type XClient struct {
	net.Conn
	HandlerClient
	sync.RWMutex
	handler  any
	tryCount uint8
	buffsize uint16
	seq      uint32
	timeout  time.Duration
	desc     grpc.ServiceDesc
	pending  map[uint32]chan Message
}

func NewClient(c net.Conn, seq uint32) Client {
	x := &XClient{
		Conn:     c,
		tryCount: 3,
		seq:      seq,
		buffsize: 16 * 1024,
		desc:     HandlerDesc,
		timeout:  time.Second * 240,
		pending:  make(map[uint32]chan Message),
	}
	x.HandlerClient = pb.NewHandlerClient(x)
	return x
}

func (x *XClient) Go(ctx context.Context, method string, req proto.Message) error {
	pusher := &Pusher{
		Conn:    x.Conn,
		timeout: x.timeout,
		seq:     math.MaxUint32,
		desc:    x.desc,
		method:  filepath.Base(method),
	}
	if err := pusher.Push(context.Background(), req); err != nil {
		return err
	}
	return nil
}

func (x *XClient) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	for i := 0; i < int(x.tryCount); i++ {
		ch, seq, err := x.newCaller()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		pusher := &Pusher{
			Conn:    x.Conn,
			timeout: x.timeout,
			seq:     seq,
			desc:    x.desc,
			method:  filepath.Base(method),
		}
		if err := pusher.Push(context.Background(), args.(proto.Message)); err != nil {
			x.done(seq)
			return err
		}
		select {
		case <-ctx.Done():
			x.done(seq)
			return errors.New("timeout")
		case message := <-ch:
			close(ch)
			if err := proto.Unmarshal(message.payload(), reply.(proto.Message)); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (x *XClient) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}

func (x *XClient) Register(a any) error {
	if x.handler != nil {
		return errors.New("x.svr  != nil")
	}
	x.handler = a
	go x.pull(context.Background())
	return nil
}

func (x *XClient) Keeplive(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		if _, err := x.HandlerClient.Ping(ctx, &pb.PingRequest{Message: []byte("ping")}); err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (x *XClient) Close() error {
	return x.Conn.Close()
}

func (x *XClient) callback(ctx context.Context, iMessage Message) error {
	ch := x.done(iMessage.ack())
	if ch == nil {
		return fmt.Errorf("%d ch == nil", iMessage.ack())
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	select {
	case ch <- iMessage:
		return nil
	case <-timeoutCtx.Done():
		close(ch)
		return nil
	}
}

func (x *XClient) pull(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
		if err := x.Close(); err != nil {
			log.Error(err.Error())
		}
	}()
	buf := bufio.NewReaderSize(x.Conn, int(x.buffsize))
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.Conn.SetReadDeadline(time.Now().Add(x.timeout)); err != nil {
			return err
		}
		message, err := decode(buf)
		if err != nil {
			return err
		}
		method := message.method()
		if method >= uint16(len(x.desc.Methods)) {
			return errors.New("kind >= len(x.desc.Methods)")
		}
		if message.ack() > 0 {
			return x.callback(ctx, message)
		}
		dec := func(in any) error {
			if err := proto.Unmarshal(message.payload(), in.(proto.Message)); err != nil {
				return err
			}
			return nil
		}
		iResponse, err := x.desc.Methods[method].Handler(x.handler, ctx, dec, nil)
		if err != nil {
			return err
		}
		if message.seq() == math.MaxUint32 {
			return nil
		}
		b, err := encode(0, message.seq(), 0, iResponse.(proto.Message))
		if err != nil {
			return err
		}
		if _, err := x.Conn.Write(b); err != nil {
			return err
		}
	}
}

func (x *XClient) newCaller() (chan Message, uint32, error) {
	x.Lock()
	defer x.Unlock()
	seq := x.seq + 1
	if _, ok := x.pending[seq]; ok {
		return nil, 0, errors.New("ok")
	}
	done := make(chan Message, 1)
	x.pending[seq] = done
	switch seq {
	case math.MaxUint32 - math.MaxUint16:
		x.seq = math.MaxUint32 / 2
	case math.MaxUint32/2 - math.MaxUint16:
		x.seq = 0
	default:
		x.seq = seq
	}
	return done, seq, nil
}

func (x *XClient) done(seq uint32) chan Message {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[seq]; ok {
		delete(x.pending, seq)
		return v
	}
	return nil
}
