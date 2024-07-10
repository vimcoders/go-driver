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

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Client interface {
	grpc.ClientConnInterface
	Register(any) error
	Keeplive(context.Context, proto.Message) error
	Go(context.Context, string, proto.Message) error
	RemoteAddr() net.Addr
	Close() error
}

type XClient struct {
	net.Conn
	grpc.ClientConnInterface
	sync.RWMutex
	handler  any
	buffsize uint16
	seq      uint32
	timeout  time.Duration
	desc     grpc.ServiceDesc
	pending  map[uint32]*stream
	streams  *sync.Pool
}

func NewClient(c net.Conn, desc grpc.ServiceDesc) Client {
	return newClient(c, desc)
}

func newClient(c net.Conn, desc grpc.ServiceDesc) Client {
	x := &XClient{
		Conn:     c,
		buffsize: 8 * 1024,
		desc:     desc,
		timeout:  time.Second * 240,
		pending:  make(map[uint32]*stream),
	}
	x.streams = &sync.Pool{
		New: func() any {
			seq := x.seq + 1
			x.seq = seq % math.MaxUint32
			return &stream{
				seq:     seq,
				Conn:    x.Conn,
				signal:  make(chan Message, 1),
				timeout: x.timeout,
			}
		},
	}
	return x
}

func (x *XClient) Keeplive(ctx context.Context, ping proto.Message) error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		case <-ticker.C:
		}
		if err := x.Go(ctx, "Ping", ping); err != nil {
			log.Error(err.Error())
			return err
		}
	}
}

func (x *XClient) Close() error {
	return x.Conn.Close()
}

func (x *XClient) Register(a any) error {
	if x.handler != nil {
		return errors.New("x.svr  != nil")
	}
	x.handler = a
	go x.pull(context.Background())
	return nil
}

func (x *XClient) Go(ctx context.Context, method string, req proto.Message) error {
	for methodId := 0; methodId < len(x.desc.Methods); methodId++ {
		if filepath.Base(method) != x.desc.Methods[methodId].MethodName {
			continue
		}
		if err := x.push(0, uint16(methodId), req); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("%s not registed", method)
}

func (x *XClient) Invoke(ctx context.Context, methodName string, req any, reply any, opts ...grpc.CallOption) (err error) {
	for method := 0; method < len(x.desc.Methods); method++ {
		if x.desc.Methods[method].MethodName != filepath.Base(methodName) {
			continue
		}
		if err := x.invoke(ctx, uint16(method), req.(proto.Message), reply.(proto.Message)); err != nil {
			return err
		}
		return nil
	}
	return errors.New(methodName)
}

func (x *XClient) invoke(ctx context.Context, method uint16, req, reply proto.Message) (err error) {
	stream := x.streams.Get().(*stream)
	x.addTask(stream)
	response, err := stream.push(ctx, method, req)
	if err != nil {
		x.done(stream.seq)
		return err
	}
	defer func() {
		response.reset()
		x.streams.Put(stream)
	}()
	if err := proto.Unmarshal(response.payload(), reply); err != nil {
		return err
	}
	return nil
}

func (x *XClient) push(seq uint32, method uint16, req proto.Message) (err error) {
	buf, err := encode(seq, method, req)
	if err != nil {
		return err
	}
	if err := x.SetWriteDeadline(time.Now().Add(x.timeout)); err != nil {
		return err
	}
	if _, err := buf.WriteTo(x.Conn); err != nil {
		return err
	}
	return nil
}

func (x *XClient) pull(ctx context.Context) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
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
		if err := x.handle(ctx, iMessage); err != nil {
			return fmt.Errorf("%v %s", iMessage, err.Error())
		}
	}
}

func (x *XClient) handle(ctx context.Context, iMessage Message) error {
	method, seq, payload := iMessage.method(), iMessage.seq(), iMessage.payload()
	if int(method) >= len(x.desc.Methods) {
		ch := x.done(seq)
		if ch == nil {
			return nil
		}
		clone, err := iMessage.clone()
		if err != nil {
			return err
		}
		if err := ch.invoke(clone); err != nil {
			return err
		}
		return nil
	}
	dec := func(in any) error {
		if err := proto.Unmarshal(payload, in.(proto.Message)); err != nil {
			return err
		}
		return nil
	}
	reply, err := x.desc.Methods[method].Handler(x.handler, ctx, dec, nil)
	if err != nil {
		return err
	}
	if seq > 0 {
		if err := x.push(seq, math.MaxUint16, reply.(proto.Message)); err != nil {
			return err
		}
	}
	return nil
}

func (x *XClient) addTask(s *stream) {
	x.Lock()
	defer x.Unlock()
	x.pending[s.seq] = s
}

func (x *XClient) done(seq uint32) *stream {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[seq]; ok {
		delete(x.pending, seq)
		return v
	}
	return nil
}
