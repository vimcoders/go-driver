package telegram

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go-driver/log"
)

const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/sendMessage"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

//var client *Client

type Client struct {
	Token  string `json:"token"`
	ChatId string
}

func Connect(chatId, token string) error {
	// client = &Client{
	// 	ChatId: chatId,
	// 	Token:  token,
	// }
	return nil
}

func (x *Client) Push(txt string) error {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
		}
	}()
	u, err := url.Parse(fmt.Sprintf(APIEndpoint, x.Token))
	if err != nil {
		return err
	}
	q := u.Query()
	q.Add("chat_id", x.ChatId)
	q.Add("text", txt)
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return err
	}
	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.ReadAll(resp.Body); err != nil {
		return err
	}
	return nil
}
