package handler

import (
	"context"
	"go-driver/app/proxy/driver"
	"go-driver/sqlx"
)

func (x *Handler) Connect(ctx context.Context) error {
	client, err := sqlx.Dial(x.Mysql.Host)
	if err != nil {
		panic(err.Error())
	}
	if err := client.Register(&driver.Account{}); err != nil {
		panic(err)
	}
	x.Client = client
	return nil
}
