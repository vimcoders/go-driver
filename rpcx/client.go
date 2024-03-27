package rpcx

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/vimcoders/go-driver/pb"

	"github.com/vimcoders/go-driver/message"

	"github.com/vimcoders/go-driver/log"

	"github.com/vimcoders/go-driver/driver"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

const (
	ReaderBuffsize = 16 * 1024
	WriterBuffsize = 16 * 1024
	Header         = 4
	Timeout        = time.Second * 120
)

type Connect struct {
	net.Conn
	OnMessage func(request proto.Message) (proto.Message, error)
	driver.Unmarshaler
	driver.Marshaler
	Closed  func()
	Timeout time.Duration
}

func (x *Connect) Read(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
		if err != nil {
			log.Error(err.Error(), x.RemoteAddr().String())
			debug.PrintStack()
		}
		x.Close()
	}()
	buffer := bufio.NewReaderSize(x.Conn, ReaderBuffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		headerBytes, err := buffer.Peek(Header)
		if err != nil {
			return err
		}
		header := binary.BigEndian.Uint32(headerBytes)
		if int(header) > buffer.Size() {
			return fmt.Errorf("header %v too long %v %v", header, headerBytes, buffer)
		}
		bodyBytes, err := buffer.Peek(int(header) + len(headerBytes))
		if err != nil {
			return err
		}
		var message pb.Message
		if err := proto.Unmarshal(bodyBytes[Header:], &message); err != nil {
			return err
		}
		if _, err := buffer.Discard(len(bodyBytes)); err != nil {
			return err
		}
		request, err := x.Unmarshal(message.Body)
		if err != nil {
			return err
		}
		response, err := x.OnMessage(request)
		if err != nil {
			log.Error(err.Error())
		}
		x.Push(context.Background(), message.RequestId, response)
	}
}

func (x *Connect) Push(ctx context.Context, requstId string, push proto.Message) (err error) {
	b, err := x.Marshal(push)
	if err != nil {
		return err
	}
	message := &pb.Message{
		Body:      b,
		RequestId: requstId,
	}
	response, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	if err := x.SetWriteDeadline(time.Now().Add(time.Second * 120)); err != nil {
		return err
	}
	buffer := make(driver.Buffer, 0, len(response)+4)
	if err := binary.Write(&buffer, binary.BigEndian, uint32(len(response))); err != nil {
		return err
	}
	buffer.Write(response)
	if _, err := x.Conn.Write(buffer); err != nil {
		return err
	}
	return nil
}

type Client struct {
	net.Conn
	sync.RWMutex
	pending   map[string]chan proto.Message
	OnMessage func(request *pb.Message) (proto.Message, error)
	OnPush    func(push proto.Message)
	driver.Unmarshaler
	driver.Marshaler
	Closed  func()
	Timeout time.Duration
}

func (x *Client) Read(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
		if err != nil {
			log.Error(err.Error(), x.RemoteAddr().String())
			debug.PrintStack()
		}
		x.Close()
	}()
	buffer := bufio.NewReaderSize(x.Conn, ReaderBuffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		if err := x.SetReadDeadline(time.Now().Add(x.Timeout)); err != nil {
			return err
		}
		headerBytes, err := buffer.Peek(Header)
		if err != nil {
			return err
		}
		header := binary.BigEndian.Uint32(headerBytes)
		if int(header) > buffer.Size() {
			return fmt.Errorf("header %v too long", header)
		}
		bodyBytes, err := buffer.Peek(int(header) + len(headerBytes))
		if err != nil {
			return err
		}
		var message pb.Message
		if err := proto.Unmarshal(bodyBytes[Header:], &message); err != nil {
			return err
		}
		x.OnMessage(&message)
		if _, err := buffer.Discard(len(bodyBytes)); err != nil {
			return err
		}
	}
}

func NewClient(c net.Conn, messages []proto.Message) *Client {
	encoder := message.NewProtobuf(messages...)
	client := &Client{
		Conn:        c,
		pending:     make(map[string]chan proto.Message, 100),
		Unmarshaler: encoder,
		Marshaler:   encoder,
		Timeout:     time.Second * 30,
	}
	client.OnMessage = func(request *pb.Message) (proto.Message, error) {
		ch := client.done(request.RequestId)
		if ch == nil {
			return nil, nil
		}
		response, err := client.Unmarshal(request.Body)
		if err != nil {
			return nil, err
		}
		ch <- response
		return nil, nil
	}
	go client.Ping(context.Background())
	go client.Read(context.Background())
	return client
}

func (x *Client) Ping(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
		if err != nil {
			log.Error(err.Error())
		}
		x.Close()
	}()
	ticker := time.NewTicker(time.Second * 60)
	for range ticker.C {
		x.Call(context.Background(), 0, &pb.LoginRequest{Token: "ping"})
	}
	return nil
}

func (x *Client) Call(ctx context.Context, id int64, args proto.Message) (reply proto.Message, err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error(), args)
		}
	}()
	b, err := x.Marshal(args)
	if err != nil {
		return nil, err
	}
	requestId := uuid.NewString()
	done := x.newPending(requestId)
	request, err := proto.Marshal(&pb.Message{Id: id, Body: b, RequestId: requestId})
	if err != nil {
		return nil, err
	}
	if err := x.SetWriteDeadline(time.Now().Add(time.Second * 120)); err != nil {
		return nil, err
	}
	buffer := make(driver.Buffer, 0, len(request)+4)
	if err := binary.Write(&buffer, binary.BigEndian, uint32(len(request))); err != nil {
		return nil, err
	}
	buffer.Write(request)
	if _, err := x.Conn.Write(buffer); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, errors.New("timeout")
	case v := <-done:
		close(done)
		return v, nil
	}
}

func (x *Client) newPending(id string) chan proto.Message {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[id]; ok {
		return v
	}
	done := make(chan proto.Message, 1)
	x.pending[id] = done
	return done
}

func (x *Client) done(id string) chan proto.Message {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[id]; ok {
		delete(x.pending, id)
		return v
	}
	return nil
}
