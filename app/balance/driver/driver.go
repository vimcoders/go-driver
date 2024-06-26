package driver

import (
	"go-driver/driver"
)

type User = driver.User
type Buffer = driver.Buffer
type Marshal = driver.Marshal
type Unmarshal = driver.Unmarshal
type ResponsePusher = driver.ResponsePusher

func NewBuffer(size int) Buffer {
	return driver.NewBuffer(size)
}
