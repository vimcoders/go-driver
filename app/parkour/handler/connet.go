package handler

import (
	"context"
)

func (x *Handler) Connect(ctx context.Context) error {
	// log.Info(x.Etcd.Endpoints)
	// cli, err := etcd.New(etcd.Config{
	// 	Endpoints:   []string{x.Etcd.Endpoints},
	// 	DialTimeout: 5 * time.Second,
	// })
	// if err != nil {
	// 	return err
	// }
	// log.Info(x.Option.Mongo.Host, x.Option.Mongo.DB)
	// mongo, err := mongox.Connect(x.Option.Mongo.Host, x.Option.Mongo.DB)
	// if err != nil {
	// 	return err
	// }
	// handler.Mongo = mongo
	// handler.Client = cli
	return nil
}
