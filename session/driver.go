// TCP，UDP 接入层
// TODO:: 熔断，限流，降级
package session

import (
	"context"
)

type Handler interface {
	Handle(ctx context.Context, request Request) error
}
