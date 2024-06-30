package rpcx

import (
	"bufio"
	"context"
	"errors"
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
	svr any
	net.Conn
	sync.RWMutex
	pending  map[uint32]chan Message
	Buffsize uint16
	Timeout  time.Duration
	seq      uint32
	desc     grpc.ServiceDesc
	HandlerClient
}

func NewClient(c net.Conn, seq uint32) Client {
	x := &XClient{
		Conn:     c,
		pending:  make(map[uint32]chan Message),
		Buffsize: 16 * 1024,
		Timeout:  time.Second * 240,
		seq:      seq,
		desc:     HandlerDesc,
	}
	x.HandlerClient = pb.NewHandlerClient(x)
	return x
}

func (x *XClient) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	for i := 0; i < math.MaxInt32; i++ {
		call, seq, err := x.newCaller()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		pusher := &Pusher{
			Conn:       x.Conn,
			timeout:    x.Timeout,
			seq:        seq,
			desc:       x.desc,
			methodName: filepath.Base(method),
		}
		if err := pusher.Push(context.Background(), args.(proto.Message)); err != nil {
			x.done(seq)
			return err
		}
		select {
		case <-ctx.Done():
			return errors.New("timeout")
		case message := <-call:
			close(call)
			if err := proto.Unmarshal(message.payload(), reply.(proto.Message)); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (x *XClient) Go(ctx context.Context, method string, req proto.Message) error {
	pusher := &Pusher{
		Conn:       x.Conn,
		timeout:    x.Timeout,
		seq:        math.MaxUint32,
		desc:       x.desc,
		methodName: filepath.Base(method),
	}
	if err := pusher.Push(context.Background(), req); err != nil {
		return err
	}
	return nil
}

func (x *XClient) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}

func (x *XClient) Register(svr any) error {
	if x.svr != nil {
		return errors.New("x.svr  != nil")
	}
	x.svr = svr
	go x.pull(context.Background())
	return nil
}

func (x *XClient) Keeplive(ctx context.Context) error {
	// for {
	// 	if err := x.Go(ctx, "Ping", &pb.PingRequest{Message: []byte("ping")}); err != nil {
	// 		log.Error(err.Error())
	// 		return err
	// 	}
	// }
	for {
		if _, err := x.Ping(ctx, &pb.PingRequest{Message: []byte("ping")}); err != nil {
			log.Error(err.Error())
			return err
		}
	}
}

func (x *XClient) Close() error {
	return x.Conn.Close()
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
	buffer := bufio.NewReaderSize(x.Conn, int(x.Buffsize))
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.Conn.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		message, err := decode(buffer)
		if err != nil {
			return err
		}
		if err := x.handleGrpc(ctx, message); err != nil {
			return err
		}
	}
}

func (x *XClient) handleGrpc(ctx context.Context, message Message) (err error) {
	kind := message.kind()
	if kind >= uint16(len(x.desc.Methods)) {
		return errors.New("kind >= len(x.desc.Methods)")
	}
	seq, ack := message.seq(), message.ack()
	if ack > 0 {
		call := x.done(ack)
		if call == nil {
			return errors.New("done == nil")
		}
		timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		select {
		case call <- message:
			return nil
		case <-timeoutCtx.Done():
			return nil
		}
	}
	method := x.desc.Methods[message.kind()]
	dec := func(in any) error {
		iMessage := in.(proto.Message)
		if err := proto.Unmarshal(message.payload(), iMessage); err != nil {
			return err
		}
		return nil
	}
	iResponse, err := method.Handler(x.svr, ctx, dec, nil)
	if err != nil {
		return err
	}
	if seq == math.MaxUint32 {
		return nil
	}
	b, err := encode(0, seq, kind, iResponse.(proto.Message))
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
	for i := uint32(1); i < math.MaxInt32; i++ {
		seq := x.seq + 1
		if _, ok := x.pending[seq]; ok {
			continue
		}
		done := make(chan Message, 1)
		x.pending[seq] = done
		switch seq {
		case math.MaxUint32:
			x.seq = math.MaxUint32 / 2
		case math.MaxUint32 / 2:
			x.seq = 0
		default:
			x.seq = seq
		}
		return done, seq, nil
	}
	return nil, 0, errors.New("too many request")
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
