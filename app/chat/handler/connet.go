package handler

import (
	"context"
	"go-driver/log"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

func (x *Handler) Connect(ctx context.Context) error {
	log.Info(x.Etcd.Endpoints)
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{x.Etcd.Endpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	x.Client = cli
	return nil
}
