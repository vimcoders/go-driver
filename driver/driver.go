// 不允许调用标准库外的包，防止循环引用
package driver

import (
	"context"
	"net"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Pool[T any] struct {
	sync.RWMutex
	Size int
	buf  []T
}

func NewPool[T any](size int) *Pool[T] {
	return &Pool[T]{
		Size: size,
	}
}

func (x *Pool[T]) Get() T {
	x.Lock()
	defer x.Unlock()
	if len(x.buf) <= 0 {
		var t T
		return t
	}
	return x.buf[0]
}

func (x *Pool[T]) Put(t T) {
	x.Lock()
	defer x.Unlock()
	if len(x.buf) >= x.Size {
		return
	}
	x.buf = append(x.buf, t)
}

type Unmarshal interface {
	Unmarshal(b []byte) (proto.Message, error)
}

type Marshal interface {
	Marshal(i proto.Message) ([]byte, error)
}

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

type ResponsePusher interface {
	Push(context.Context, proto.Message) error
}

type Buffer []byte

func NewBuffer(size int) Buffer {
	return make([]byte, size)
}

func (b *Buffer) Reset() {
	*b = (*b)[:0]
}

func (b *Buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *Buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}
