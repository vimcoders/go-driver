// 不允许调用标准库外的包，防止循环引用
package driver

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	Write(i any) (int, error)
}

type ResponseWriter struct {
	W       http.ResponseWriter `json:"-"`
	Code    int32               `json:"code"`
	Message string              `json:"message"`
	Data    any                 `json:"data"`
}

func (x *ResponseWriter) Write(i any) (int, error) {
	x.Data = i
	x.Code = http.StatusOK
	b, err := json.Marshal(x)
	if err != nil {
		return 0, err
	}
	return x.W.Write(b)
}
