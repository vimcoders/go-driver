package benchmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-driver/driver"
	"go-driver/log"
	"go-driver/pb"
	"go-driver/tcp"
	"io"
	"math/rand"
	"net"
	"net/http"
	"sync"

	"google.golang.org/grpc"
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
	Token    string
	tcp.Client
	grpc.ServiceDesc
	sync.Pool
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

func (x *Client) Ping(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
		default:
			b, err := proto.Marshal(&pb.PingRequest{})
			if err != nil {
				panic(err)
			}
			response := x.Pool.Get().(*driver.Message)
			response.WriteUint16(uint16(4 + len(b)))
			response.WriteUint16(0)
			response.Write(b)
			if _, err := response.WriteTo(x.Client); err != nil {
				panic(err)
			}
			response.Reset()
			x.Pool.Put(response)
		}
	}
}

func (x *Client) ServeTCP(ctx context.Context, buf []byte) error {
	return nil
}

func (x *Client) ServeKCP(ctx context.Context, buf []byte) error {
	return x.Handle(ctx, buf)
}

func (x *Client) ServeQUIC(ctx context.Context, buf []byte) error {
	return x.Handle(ctx, buf)
}

func (x *Client) Handle(ctx context.Context, buf []byte) error {
	// var request driver.Message = buf
	// method, payload := request.Method(), request.Payload()
	// dec := func(in any) error {
	// 	if err := proto.Unmarshal(payload, in.(proto.Message)); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }
	// _, err := x.Methods[method].Handler(x, ctx, dec, nil)
	// if err != nil {
	// 	log.Error(err.Error())
	// }
	return nil
}
