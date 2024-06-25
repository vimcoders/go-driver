// TCP，UDP 接入层
// TODO:: 熔断，限流，降级
package handle

import (
	"context"
)

// 我们将在这里定义一个接口来处理我们解析出来的二进制流
type Handler interface {
	Handle(ctx context.Context, request Request) error
}
