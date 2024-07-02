package grpcx

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"go-driver/log"
	"go-driver/pb"
	"go-driver/quicx"

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

func Dial(network string, addr string) (Client, error) {
	switch network {
	case "udp":
		conn, err := quicx.Dial(addr, &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"quic-echo-example"},
			MaxVersion:         tls.VersionTLS13,
		}, &quicx.Config{
			MaxIdleTimeout: time.Minute,
		})
		if err != nil {
			return nil, err
		}
		return newClient(conn, 0), nil
	case "tcp":
		fallthrough
	case "tcp4":
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		return newClient(conn, 0), nil
	}
	return nil, fmt.Errorf("%s unkonw", network)
}

func NewClient(c net.Conn) Client {
	return newClient(c, math.MaxUint32/2)
}

func newClient(c net.Conn, seq uint32) Client {
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

func (x *XClient) Go(ctx context.Context, req proto.Message) error {
	pusher := &Pusher{
		Conn:    x.Conn,
		timeout: x.timeout,
		seq:     math.MaxUint32,
		desc:    x.desc,
	}
	if err := pusher.Push(context.Background(), req); err != nil {
		return err
	}
	return nil
}

func (x *XClient) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) (err error) {
	for i := 0; i < int(x.tryCount); i++ {
		if err := x.invoke(ctx, method, args, reply); err != nil {
			log.Error(err.Error(), method, args, reply)
			continue
		}
		return nil
	}
	return errors.New("try many invoke")
}

func (x *XClient) invoke(ctx context.Context, _ string, args any, reply any) (err error) {
	ch, seq, err := x.newCaller()
	if err != nil {
		return err
	}
	pusher := &Pusher{
		Conn:    x.Conn,
		timeout: x.timeout,
		seq:     seq,
		desc:    x.desc,
	}
	if err := pusher.Push(context.Background(), args.(proto.Message)); err != nil {
		x.done(seq)
		return err
	}
	select {
	case <-ctx.Done():
		x.done(seq)
		log.Error("invoke cancel")
	case iMessage := <-ch:
		if err := proto.Unmarshal(iMessage.payload(), reply.(proto.Message)); err != nil {
			log.Error(iMessage, seq)
			return err
		}
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
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		if _, err := x.Ping(ctx, &pb.PingRequest{Message: []byte("ping")}); err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (x *XClient) Close() error {
	return x.Conn.Close()
}

func (x *XClient) callback(ctx context.Context, ack uint32, clone Message) error {
	ch := x.done(ack)
	if ch == nil {
		log.Errorf("%d ch == nil", ack)
		return nil
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	select {
	case ch <- clone:
	case <-timeoutCtx.Done():
		log.Error("<-timeoutCtx.Done()")
	}
	return nil
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
		iMessage, err := decode(buf)
		if err != nil {
			return err
		}
		if err := x.handle(ctx, iMessage.clone()); err != nil {
			return err
		}
	}
}

func (x *XClient) handle(ctx context.Context, clone Message) error {
	method := clone.method()
	if method >= uint16(len(x.desc.Methods)) {
		return errors.New("kind >= len(x.desc.Methods)")
	}
	seq, ack, payload := clone.seq(), clone.ack(), clone.payload()
	if ack > 0 {
		return x.callback(ctx, ack, clone)
	}
	dec := func(in any) error {
		if err := proto.Unmarshal(payload, in.(proto.Message)); err != nil {
			return err
		}
		return nil
	}
	iResponse, err := x.desc.Methods[method].Handler(x.handler, ctx, dec, nil)
	if err != nil {
		return err
	}
	if seq == math.MaxUint32 {
		return nil
	}
	b, err := encode(0, seq, 0, iResponse.(proto.Message))
	if err != nil {
		return err
	}
	if _, err := x.Conn.Write(b); err != nil {
		return err
	}
	return nil
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
