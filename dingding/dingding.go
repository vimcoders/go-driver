package dingding

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client *Client

type Client struct {
	AccessToken string
	Secret      string
	EnableAt    bool
	AtAll       bool
}

func Connect(token, secret string) error {
	client = &Client{
		AccessToken: token,
		Secret:      secret,
		EnableAt:    true,
		AtAll:       true,
	}
	return nil
}

func Push(txt string) error {
	return client.Push(txt)
}

// SendMessage Function to send message
//
//goland:noinspection GoUnhandledErrorResult
func (x *Client) Push(s string, at ...string) error {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": s,
		},
	}
	if x.EnableAt {
		if x.AtAll {
			if len(at) > 0 {
				return errors.New("the parameter \"AtAll\" is \"true\", but the \"at\" parameter of SendMessage is not empty")
			}
			msg["at"] = map[string]interface{}{
				"isAtAll": x.AtAll,
			}
		} else {
			msg["at"] = map[string]interface{}{
				"atMobiles": at,
				"isAtAll":   x.AtAll,
			}
		}
	} else {
		if len(at) > 0 {
			return errors.New("the parameter \"EnableAt\" is \"false\", but the \"at\" parameter of SendMessage is not empty")
		}
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	wh := "https://oapi.dingtalk.com/robot/send?access_token=" + x.AccessToken
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, x.Secret)
	sign := x.hmacSha256(stringToSign, x.Secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", wh, timestamp, sign)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.ReadAll(resp.Body); err != nil {
		return err
	}
	return nil
}

func (x *Client) hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
