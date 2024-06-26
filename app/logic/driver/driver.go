// 应该在进程内自己实现的数据结构，为了提高复用性，从外部引入的
package driver

import (
	"go-driver/driver"
)

type User = driver.User
type ResponsePusher = driver.ResponsePusher
