// 不允许调用标准库外的包，防止循环引用
package driver

import (
	"context"
	"net"

	"google.golang.org/protobuf/proto"
)

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
