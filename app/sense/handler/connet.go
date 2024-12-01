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
	// response, err := cli.Get(ctx, x.Etcd.Join("logic"))
	// if err != nil {
	// 	return err
	// }
	// for i := 0; i < len(response.Kvs); i++ {
	// 	var service etcdx.Service
	// 	if err := json.Unmarshal(response.Kvs[i].Value, &service); err != nil {
	// 		log.Error(err.Error())
	// 		continue
	// 	}
	// 	// client, err := grpcx.Dial("udp", service.LocalAddr)
	// 	// if err != nil {
	// 	// 	log.Error(err.Error())
	// 	// 	continue
	// 	// }
	// 	// if _, err := client.Ping(ctx, &pb.PingRequest{}); err != nil {
	// 	// 	log.Error(err.Error())
	// 	// 	continue
	// 	// }
	// 	handler.Client = cli
	// 	//handler.rpc = client
	// 	return nil
	// }
	return nil
}
