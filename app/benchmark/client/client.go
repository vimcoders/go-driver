package benchmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-driver/driver"
	"go-driver/handle"
	"go-driver/log"
	"go-driver/pb"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	"google.golang.org/protobuf/proto"
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
	h        *handle.Handle
	driver.Marshal
	driver.Unmarshal
	Token string
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

func (x *Client) Login() error {
	conn, err := net.Dial("tcp", x.CometUrl)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	// conn, err := quicx.Dial(x.CometUrl, &tls.Config{
	// 	InsecureSkipVerify: true,
	// 	NextProtos:         []string{"quic-echo-example"},
	// 	MaxVersion:         tls.VersionTLS13,
	// }, &quicx.Config{
	// 	MaxIdleTimeout: time.Minute,
	// })
	//conn, err := net.Dial("tcp", response[i].Addr)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return err
	// }
	if err := x.Register(); err != nil {
		return err
	}
	x.h = handle.NewHandle(conn)
	x.h.Handler = x
	go x.h.Pull(context.Background())
	go x.Keeplive(context.Background())
	if err := x.Push(context.Background(), &pb.LoginRequest{Token: x.Token}); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (x *Client) LoginResponse(response *pb.LoginResponse) {

}

func (x *Client) Handle(ctx context.Context, request handle.Request) error {
	// message, _, err := x.Unmarshal.Unmarshal(request)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return
	// }
	// method := reflect.ValueOf(x).MethodByName(string(proto.MessageName(message).Name()))
	// method.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(message)})
	return nil
}

func (x *Client) Keeplive(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		if err := x.Ping(ctx); err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (x *Client) Ping(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	if err := x.Push(ctx, &pb.PingRequest{}); err != nil {
		return err
	}
	return nil
}

func (x *Client) Push(ctx context.Context, message proto.Message) error {
	response, err := x.Marshal.Marshal(message)
	if err != nil {
		return err
	}
	if _, err := x.h.Push(ctx, response); err != nil {
		return err
	}
	return nil
}
