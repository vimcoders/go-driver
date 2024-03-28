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

	"github.com/google/uuid"
	"github.com/vimcoders/go-driver/driver"
	"github.com/vimcoders/go-driver/log"
	"github.com/vimcoders/go-driver/pb"
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
	OnMessage func(request *Request) (*Response, error)
	Closed    func()
	Timeout   time.Duration
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
		var request pb.Request
		if err := proto.Unmarshal(bodyBytes[Header:], &request); err != nil {
			return err
		}
		if _, err := buffer.Discard(len(bodyBytes)); err != nil {
			return err
		}
		response, err := x.OnMessage(&Request{RequestId: request.RequestId, Message: request.Message})
		if err != nil {
			log.Error(err.Error())
		}
		response.RequestId = request.RequestId
		x.Push(context.Background(), response)
	}
}

func (x *Connect) Push(ctx context.Context, response *Response) (err error) {
	iResponse, err := proto.Marshal(response.ToMessage())
	if err != nil {
		return err
	}
	if err := x.SetWriteDeadline(time.Now().Add(time.Second * 120)); err != nil {
		return err
	}
	buffer := make(driver.Buffer, 0, len(iResponse)+4)
	if err := binary.Write(&buffer, binary.BigEndian, uint32(len(iResponse))); err != nil {
		return err
	}
	buffer.Write(iResponse)
	if _, err := x.Conn.Write(buffer); err != nil {
		return err
	}
	return nil
}

type Client struct {
	net.Conn
	sync.RWMutex
	pending map[string]chan *Response
	OnPush  func(push proto.Message)
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
		var response pb.Response
		if err := proto.Unmarshal(bodyBytes[Header:], &response); err != nil {
			return err
		}
		ch := x.done(response.RequestId)
		if ch == nil {
			continue
		}
		ch <- &Response{RequestId: response.RequestId, Message: response.Message}
		if _, err := buffer.Discard(len(bodyBytes)); err != nil {
			return err
		}
	}
}

func NewClient(c net.Conn, messages []proto.Message) *Client {
	client := &Client{
		Conn:    c,
		pending: make(map[string]chan *Response),
		Timeout: time.Second * 30,
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
		b, err := proto.Marshal(&pb.LoginRequest{Token: "ping"})
		if err != nil {
			log.Error(err.Error())
			continue
		}
		x.Call(context.Background(), &Request{Message: b})
	}
	return nil
}

func (x *Client) Call(ctx context.Context, request *Request) (reply *Response, err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error(), request)
		}
	}()
	request.RequestId = uuid.NewString()
	b, err := proto.Marshal(request.ToMessage())
	if err != nil {
		return nil, err
	}
	done := x.newPending(request.RequestId)
	if err := x.SetWriteDeadline(time.Now().Add(time.Second * 120)); err != nil {
		return nil, err
	}
	buffer := make(driver.Buffer, 0, len(b)+4)
	if err := binary.Write(&buffer, binary.BigEndian, uint32(len(b))); err != nil {
		return nil, err
	}
	buffer.Write(b)
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

func (x *Client) newPending(id string) chan *Response {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[id]; ok {
		return v
	}
	done := make(chan *Response, 1)
	x.pending[id] = done
	return done
}

func (x *Client) done(id string) chan *Response {
	x.Lock()
	defer x.Unlock()
	if v, ok := x.pending[id]; ok {
		delete(x.pending, id)
		return v
	}
	return nil
}
