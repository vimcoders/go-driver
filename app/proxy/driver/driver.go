// 应该在进程内自己实现的数据结构，为了提高复用性，从外部引入的
package driver

import "go-driver/driver"

type Buffer = driver.Buffer
type Account = driver.Account
type Response = driver.Response
type ResponsePusher = driver.ResponsePusher
type PassportLoginRequest = driver.PassportLoginRequest
type PassportLoginResponse = driver.PassportLoginResponse
type Metod = driver.Metod

func NewBuffer(size int) Buffer {
	return driver.NewBuffer(size)
}
