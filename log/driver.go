package log

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

type Buffer []byte

func NewBuffer(size int) Buffer {
	return make([]byte, size)
}

func (x *Buffer) free() {
	*x = (*x)[:0]
	bufferFree.Put(x)
}

func (b *Buffer) Write(s string) {
	*b = append(*b, s...)
}

func (x *Buffer) WriteUint32(v uint32) {
	*x = binary.BigEndian.AppendUint32(*x, v)
}

func (x *Buffer) Appendln(a ...any) {
	*x = fmt.Appendln(*x, a...)
}

func newPrinter(prefix string, a ...any) *Buffer {
	buffer := bufferFree.Get().(*Buffer)
	buffer.Write(time.Now().Format("2006-01-02 15:04:05"))
	buffer.Write(prefix)
	buffer.Appendln(a...)
	return buffer
}

func newPrinterf(prefix, format string, a ...any) *Buffer {
	return newPrinter(prefix, fmt.Sprintf(format, a...))
}

var bufferFree sync.Pool = sync.Pool{
	New: func() any {
		return &Buffer{}
	},
}

type Handler interface {
	Handle(context.Context, []byte) error
	Close() error
}
