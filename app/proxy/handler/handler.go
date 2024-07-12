package handler

import (
	"go-driver/conf"
	"go-driver/driver"
	"go-driver/sqlx"
)

type Handler struct {
	Opt *conf.Conf
	sqlx.Client
}

func MakeHandler(opt *conf.Conf) *Handler {
	client, err := sqlx.Dial(opt.Mysql.Host)
	if err != nil {
		panic(err.Error())
	}
	if err := client.Register(&driver.Account{}); err != nil {
		panic(err)
	}
	return &Handler{Opt: opt, Client: client}
}
