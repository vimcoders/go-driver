package log

import (
	"context"
	"go-driver/driver"
)

type Buffer = driver.Buffer

func NewBuffer(size int) Buffer {
	return driver.NewBuffer(size)
}

type Handler interface {
	Handle(context.Context, []byte) error
	Close() error
}
