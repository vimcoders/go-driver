package benchmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-driver/driver"
	"go-driver/log"
	"go-driver/tcp"
	"io"
	"math/rand"
	"net"
	"net/http"
)

type ResponseWriter[T any] struct {
	W       http.ResponseWriter `json:"-"`
	Code    int32               `json:"code"`
	Message string              `json:"message"`
	Data    T                   `json:"data"`
}

type Client struct {
	Url      string
	CometUrl string
	Token    string
	tcp.Client
	//messages []proto.Message
}

func (x *Client) Register() error {
	b, err := json.Marshal(&driver.PassportLoginRequest{
		Passport: fmt.Sprintf("%d", rand.Int63()),
		Pwd:      fmt.Sprintf("%d", rand.Int63()),
	})
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", x.Url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	responseWriter := ResponseWriter[driver.PassportLoginResponse]{}
	if err := json.Unmarshal(body, &responseWriter); err != nil {
		return err
	}
	x.Token = responseWriter.Data.Token
	return nil
}

func (x *Client) ServeTCP(ctx context.Context, reply []byte) error {
	return nil
}

func (x *Client) Login() error {
	// conn, err := quicx.Dial(x.CometUrl, &tls.Config{
	// 	InsecureSkipVerify: true,
	// 	NextProtos:         []string{"quic-echo-example"},
	// 	MaxVersion:         tls.VersionTLS13,
	// }, &quicx.Config{
	// 	MaxIdleTimeout: time.Minute,
	// })
	conn, err := net.Dial("tcp", x.CometUrl)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	x.Client = tcp.NewClient(conn, tcp.Option{})
	if err := x.Register(); err != nil {
		return err
	}
	if err := x.Client.Register(x); err != nil {
		return err
	}
	// if err := x.Go(context.Background(), &pb.LoginRequest{Token: x.Token}); err != nil {
	// 	return err
	// }
	//go x.Keeplive(context.Background(), &pb.PingRequest{})
	return nil
}
