package driver

import (
	"context"
	"io"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type OptionFunc func(opt *Option)

type Option struct {
	Unmarshaler
	Marshaler
	OnMessage func(c Conn, v proto.Message) (proto.Message, error)
	OnPush    func(c Conn, pushId int64, v proto.Message) (proto.Message, error)
	Closed    func()
	Timeout   time.Duration
	// ReaderBuffsize int
	// WriterBuffsize int
	// Header         int
}

type Unmarshaler interface {
	Unmarshal(data []byte) (proto.Message, error)
}

type Marshaler interface {
	Marshal(v proto.Message) ([]byte, error)
}

type Conn interface {
	Read(ctx context.Context) error
	Write(ctx context.Context, v proto.Message) error
	RemoteAddr() net.Addr
	Options(opts ...OptionFunc)
	io.Closer
	SetReadDeadline(t time.Time) error
}

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

type Document interface {
	DocumentId() string
	DocumentName() string
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

type Encoder interface {
	Encode(origData []byte) ([]byte, error)
}

type Decoder interface {
	Decode(crypted []byte) ([]byte, error)
}
