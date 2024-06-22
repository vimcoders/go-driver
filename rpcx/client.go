package rpcx

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"go-driver/driver"
	"go-driver/log"
	"go-driver/pb"

	"google.golang.org/protobuf/proto"
)

const (
	ReaderBuffsize = 16 * 1024
	WriterBuffsize = 16 * 1024
	LENGTH         = 4
	TIMEOUT        = time.Second * 120
	MESSAGEID      = "message_id"
	MESSAGENAME    = "message_name"
)

type Client struct {
	net.Conn
	sync.RWMutex
	messageId int32
	pending   map[int32]chan *pb.Message
}

func NewClient(c net.Conn) *Client {
	client := &Client{
		Conn:    c,
		pending: make(map[int32]chan *pb.Message),
	}
	go client.Poll(context.Background())
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

func (x *Client) Poll(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
			debug.PrintStack()
		}
		if err := x.Close(); err != nil {
			log.Error(err.Error())
		}
	}()
	buffer := bufio.NewReaderSize(x.Conn, ReaderBuffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.SetReadDeadline(time.Now().Add(TIMEOUT)); err != nil {
			return err
		}
		bytes, err := buffer.Peek(LENGTH)
		if err != nil {
			return err
		}
		length := binary.BigEndian.Uint32(bytes)
		if int(length) > buffer.Size() {
			return fmt.Errorf("header %v too long", length)
		}
		message, err := buffer.Peek(int(length) + len(bytes))
		if err != nil {
			return err
		}
		var response pb.Message
		if err := proto.Unmarshal(message[LENGTH:], &response); err != nil {
			return err
		}
		if _, err := buffer.Discard(len(message)); err != nil {
			return err
		}
		var opt Option = response.Option
		messageId, err := strconv.Atoi(opt.Get(MESSAGEID))
		if err != nil {
			return err
		}
		done := x.done(int32(messageId))
		if done == nil {
			continue
		}
		done <- &response
	}
}

func (x *Client) Call(ctx context.Context, request proto.Message, reply proto.Message, opt ...*pb.Option) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	done, messageId := x.newPending()
	if done == nil {
		return errors.New("done == nil")
	}
	m, err := proto.Marshal(request)
	if err != nil {
		return err
	}
	message := &pb.Message{Message: m}
	message.Option = append(message.Option, opt...)
	message.Option = append(message.Option, &pb.Option{Key: MESSAGEID, Value: fmt.Sprintf("%v", messageId)})
	message.Option = append(message.Option, &pb.Option{Key: MESSAGENAME, Value: string(proto.MessageName(request).Name())})
	b, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	buffer := make(driver.Buffer, 4)
	binary.BigEndian.PutUint32(buffer, uint32(len(b)))
	buffer.Write(b)
	if err := x.SetWriteDeadline(time.Now().Add(TIMEOUT)); err != nil {
		return err
	}
	if _, err := x.Conn.Write(buffer); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		x.done(messageId)
		return errors.New("timeout")
	case v := <-done:
		close(done)
		return proto.Unmarshal(v.Message, reply)
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

func (x *Client) newPending() (chan *pb.Message, int32) {
	x.Lock()
	defer x.Unlock()
	for i := int32(1); i < math.MaxInt32; i++ {
		messageId := x.messageId + i
		if _, ok := x.pending[messageId]; ok {
			continue
		}
		done := make(chan *pb.Message, 1)
		x.pending[messageId] = done
		x.messageId = x.messageId%math.MaxInt32 + 1
		return done, messageId
	}
	return nil, 0
}

func (x *Client) done(messageId int32) chan *pb.Message {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[messageId]; ok {
		delete(x.pending, messageId)
		return v
	}
	return nil
}
