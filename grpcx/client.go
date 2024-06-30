package grpcx

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
	net.Conn
	HandlerClient
	sync.RWMutex
	pending  map[uint32]chan Message
	buffsize uint16
	timeout  time.Duration
	seq      uint32
	desc     grpc.ServiceDesc
	svr      any
}

func NewClient(c net.Conn, seq uint32) Client {
	x := &XClient{
		Conn:     c,
		pending:  make(map[uint32]chan Message),
		buffsize: 16 * 1024,
		timeout:  time.Second * 240,
		seq:      seq,
		desc:     HandlerDesc,
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
	for i := 0; i < math.MaxInt32; i++ {
		caller, seq, err := x.newCaller()
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
		case message := <-caller:
			close(caller)
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

func (x *XClient) Register(svr any) error {
	if x.svr != nil {
		return errors.New("x.svr  != nil")
	}
	x.svr = svr
	go x.pull(context.Background())
	return nil
}

func (x *XClient) Keeplive(ctx context.Context) error {
	for {
		if _, err := x.Ping(ctx, &pb.PingRequest{Message: []byte("ping")}); err != nil {
			log.Error(err.Error())
			return err
		}
	}
}

func (x *XClient) Ping(ctx context.Context, args *pb.PingRequest, opts ...grpc.CallOption) (*pb.PingResponse, error) {
	//x.Go(ctx, "Ping", &pb.PingRequest{Message: []byte("ping")})
	return x.HandlerClient.Ping(ctx, args, opts...)
}

func (x *XClient) Close() error {
	return x.Conn.Close()
}

func (x *XClient) callback(ctx context.Context, iMessage Message) error {
	caller := x.done(iMessage.ack())
	if caller == nil {
		return errors.New("done == nil")
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	select {
	case caller <- iMessage:
		return nil
	case <-timeoutCtx.Done():
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
	buffer := bufio.NewReaderSize(x.Conn, int(x.buffsize))
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.Conn.SetReadDeadline(time.Now().Add(x.timeout)); err != nil {
			return err
		}
		message, err := decode(buffer)
		if err != nil {
			return err
		}
		if err := x.handle(ctx, message); err != nil {
			return err
		}
	}
}

func (x *XClient) handle(ctx context.Context, message Message) (err error) {
	methodId := message.id()
	if methodId >= uint16(len(x.desc.Methods)) {
		return errors.New("kind >= len(x.desc.Methods)")
	}
	seq, ack := message.seq(), message.ack()
	if ack > 0 {
		return x.callback(ctx, message)
	}
	method := x.desc.Methods[methodId]
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
