package handler

import (
	"go-driver/conf"
	"go-driver/driver"
	"go-driver/log"
	"go-driver/sqlx"
)

type Handler struct {
	Opt *conf.Conf
	*sqlx.Client
}

func MakeHandler(opt *conf.Conf) *Handler {
	log.Debug(opt.Mysql.Host)
	client, err := sqlx.Connect(opt.Mysql.Host)
	if err != nil {
		panic(err.Error())
	}
	client.Register(&driver.Account{})
	return &Handler{Opt: opt, Client: client}
}
