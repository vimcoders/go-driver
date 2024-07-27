package handler

import (
	"go-driver/app/proxy/driver"
	"go-driver/sqlx"
)

type Handler struct {
	*driver.Option
	sqlx.Client
}

func MakeHandler(opt *driver.Option) *Handler {
	client, err := sqlx.Dial(opt.Mysql.Host)
	if err != nil {
		panic(err.Error())
	}
	if err := client.Register(&driver.Account{}); err != nil {
		panic(err)
	}
	return &Handler{Option: opt, Client: client}
}
