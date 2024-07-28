package handler

import (
	"context"
	"go-driver/app/proxy/driver"
	"go-driver/sqlx"
	"net/http"
)

type Handler struct {
	Option
	sqlx.Client
	trees map[string]func(w driver.Response, r *http.Request)
}

func MakeHandler(ctx context.Context) *Handler {
	opt := ParseOption()
	client, err := sqlx.Dial(opt.Mysql.Host)
	if err != nil {
		panic(err.Error())
	}
	if err := client.Register(&driver.Account{}); err != nil {
		panic(err)
	}
	return &Handler{Option: opt, Client: client}
}
