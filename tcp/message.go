package tcp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

type Message []byte

// 数据流解密
func decode(b *bufio.Reader) (Message, error) {
	headerBytes, err := b.Peek(4)
	if err != nil {
		return nil, err
	}
	header := binary.BigEndian.Uint32(headerBytes)
	if int(header) > b.Size() {
		return nil, fmt.Errorf("header %v too long", header)
	}
	request, err := b.Peek(int(header))
	if err != nil {
		return nil, err
	}
	return request, nil
}

// 数据流加密
func encode(kind uint16, message proto.Message) (Message, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var iMessage Message
	iMessage.WriteUint32(uint32(6 + len(b)))
	iMessage.WriteUint16(kind)
	if _, err := iMessage.Write(b); err != nil {
		return nil, err
	}
	return iMessage, nil
}

// 从数据流中获取协议号
func (x Message) kind() uint16 {
	return binary.BigEndian.Uint16(x[4:])
}

// 从数据流中获取包体
func (x Message) message() []byte {
	return x[6:]
}

func (x *Message) reset() {
	if cap(*x) <= 0 {
		return
	}
	*x = (*x)[:0]
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
