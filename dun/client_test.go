package dun_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/vimcoders/go-driver/dun"
)

func TestCheckText(t *testing.T) {
	params := url.Values{
		"dataId":  []string{"ebfcad1c-dba1-490c-b4de-e784c2691768"},
		"content": []string{"易盾你！"},
		//"dataType": []string{"1"},
		//"ip": []string{"123.115.77.137"},
		//"account": []string{"golang@163.com"},
		//"deviceType": []string{"4"},
		//"deviceId": []string{"92B1E5AA-4C3D-4565-A8C2-86E297055088"},
		//"callback": []string{"ebfcad1c-dba1-490c-b4de-e784c2691768"},
		//"publishTime": []string{"1479677336255"},
		//"callbackUrl": []string{"http://***"},	//主动回调地址url,如果设置了则走主动回调逻辑
	}
	client := dun.NewClient()
	fmt.Println(client.CheckText(params))
}
