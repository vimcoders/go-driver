package grpcx

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Message []byte

var pool sync.Pool = sync.Pool{
	New: func() any {
		return &Message{}
	},
}

func decode(b *bufio.Reader) (Message, error) {
	headerBytes, err := b.Peek(2)
	if err != nil {
		return nil, err
	}
	length := int(binary.BigEndian.Uint16(headerBytes))
	if length > b.Size() {
		return nil, fmt.Errorf("header %v too long", length)
	}
	iMessage, err := b.Peek(length)
	if err != nil {
		return nil, err
	}
	if _, err := b.Discard(len(iMessage)); err != nil {
		return nil, err
	}
	return iMessage, nil
}

func encode(seq uint32, method uint16, iMessage proto.Message) (Message, error) {
	b, err := proto.Marshal(iMessage)
	if err != nil {
		return nil, err
	}
	buf := pool.Get().(*Message)
	buf.WriteUint16(uint16(8 + len(b)))
	buf.WriteUint32(seq)
	buf.WriteUint16(method)
	if _, err := buf.Write(b); err != nil {
		return nil, err
	}
	return *buf, nil
}

func (x Message) length() uint16 {
	return binary.BigEndian.Uint16(x)
}

func (x Message) seq() uint32 {
	return binary.BigEndian.Uint32(x[2:])
}

func (x Message) method() uint16 {
	return binary.BigEndian.Uint16(x[6:])
}

func (x Message) payload() []byte {
	return x[8:x.length()]
}

func (x Message) clone() (Message, error) {
	clone := pool.Get().(*Message)
	if _, err := clone.Write(x); err != nil {
		return nil, err
	}
	return *clone, nil
}

func (x *Message) reset() {
	if cap(*x) <= 0 {
		return
	}
	*x = (*x)[:0]
	pool.Put(x)
}

func (x *Message) Write(p []byte) (int, error) {
	*x = append(*x, p...)
	return len(p), nil
}

func (x *Message) WriteUint32(v uint32) {
	*x = binary.BigEndian.AppendUint32(*x, v)
}

func (x *Message) WriteUint16(v uint16) {
	*x = binary.BigEndian.AppendUint16(*x, v)
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (x Message) WriteTo(w io.Writer) (n int64, err error) {
	if nBytes := len(x); nBytes > 0 {
		m, e := w.Write(x)
		if m > nBytes {
			panic("bytes.Buffer.WriteTo: invalid Write count")
		}
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// Buffer is now empty; reset.
	x.reset()
	return n, nil
}
