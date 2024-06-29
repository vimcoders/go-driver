package rpcx

import (
	"bufio"
	"context"
	"errors"
	"math"
	"net"
	"sync"
	"time"

	"go-driver/log"
	"go-driver/pb"

	"google.golang.org/protobuf/proto"
)

type XClient struct {
	Handler Handler
	net.Conn
	sync.RWMutex
	messageId uint32
	pending   map[uint32]chan proto.Message
	Buffsize  uint16
	Timeout   time.Duration
	messages  []proto.Message
	seq       uint16
}

func NewClient(c net.Conn, seq uint16) Client {
	x := &XClient{
		Conn:     c,
		pending:  make(map[uint32]chan proto.Message),
		Buffsize: 16 * 1024,
		Timeout:  time.Second * 240,
		messages: messages,
		seq:      seq,
	}
	return x
}

func (x *XClient) Register(h interface{}) error {
	if x.Handler != nil {
		return errors.New("x.Handler != nil")
	}
	handler, ok := h.(Handler)
	if !ok {
		return errors.New("Handler !ok")
	}
	x.Handler = handler
	go x.pull(context.Background())
	return nil
}

func (x *XClient) Keeplive(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		if err := x.Ping(ctx); err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (x *XClient) Ping(ctx context.Context) error {
	if _, err := x.Call(ctx, &pb.PingRequest{}); err != nil {
		return err
	}
	return nil
}

func (x *XClient) Call(ctx context.Context, request proto.Message) (reply proto.Message, err error) {
	for i := 0; i < math.MaxInt32; i++ {
		call, messageId, err := x.addCall()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		pusher := &Pusher{
			Conn:     x.Conn,
			timeout:  x.Timeout,
			messages: x.messages,
			seq:      messageId,
		}
		if err := pusher.push(context.Background(), request); err != nil {
			x.done(messageId)
			return nil, err
		}
		select {
		case <-ctx.Done():
			x.done(messageId)
			return nil, errors.New("timeout")
		case message := <-call:
			close(call)
			return message, nil
		}
	}
	return nil, errors.New("try many request")
}

func (x *XClient) Go(ctx context.Context, request proto.Message) (err error) {
	pusher := &Pusher{
		Conn:     x.Conn,
		timeout:  x.Timeout,
		messages: x.messages,
		seq:      math.MaxUint32,
	}
	return pusher.push(context.Background(), request)
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
		go x.handle(ctx, message)
	}
}

func (x *XClient) handle(ctx context.Context, message Message) error {
	response, err := x.new(message.kind())
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(message.message(), response); err != nil {
		return err
	}
	seq, ack := message.seq(), message.ack()
	if seq == math.MaxUint32 {
		return x.serveCast(ctx, response)
	}
	if seq > 0 {
		return x.serveCall(ctx, seq, response)
	}
	call := x.done(ack)
	if call == nil {
		return errors.New("done == nil")
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	select {
	case call <- response:
		return nil
	case <-timeoutCtx.Done():
		return nil
	}
}

func (x *XClient) new(kind uint16) (proto.Message, error) {
	if kind >= uint16(len(x.messages)) {
		return nil, errors.New("kind >= uint16(len(x.messages))")
	}
	return x.messages[kind].ProtoReflect().New().Interface(), nil
}

func (x *XClient) serveCall(ctx context.Context, seq uint32, message proto.Message) error {
	// methodName := proto.MessageName(message).Name()
	// method := reflect.ValueOf(x.Handler).MethodByName(string(methodName))
	// if ok := method.IsValid(); !ok {
	// 	return errors.New("method.IsValid(); !ok")
	// }
	// args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(message)}
	result, err := x.Handler.Call(ctx, message)
	if err != nil {
		return err
	}
	pusher := &Pusher{
		Conn:     x.Conn,
		timeout:  x.Timeout,
		messages: x.messages,
		ack:      seq,
	}
	return pusher.push(ctx, result)
}

func (x *XClient) serveCast(ctx context.Context, message proto.Message) error {
	return x.Handler.Go(ctx, message)
}

func (x *XClient) addCall() (chan proto.Message, uint32, error) {
	x.Lock()
	defer x.Unlock()
	limit := uint32(math.MaxUint16) + uint32(x.seq)
	for i := uint32(x.seq); i < limit; i++ {
		messageId := x.messageId + i
		if _, ok := x.pending[messageId]; ok {
			continue
		}
		done := make(chan proto.Message, 1)
		x.pending[messageId] = done
		x.messageId = messageId % math.MaxUint16
		return done, messageId, nil
	}
	return nil, 0, errors.New("too many request")
}

func (x *XClient) done(seq uint32) chan proto.Message {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[seq]; ok {
		delete(x.pending, seq)
		return v
	}
	return nil
}
