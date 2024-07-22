package log

import (
	"context"
	"go-driver/driver"
)

type Buffer = driver.Buffer

type Handler interface {
	Handle(context.Context, []byte) error
	Close() error
}
