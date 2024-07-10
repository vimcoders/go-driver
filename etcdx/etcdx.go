package etcdx

import (
	"context"
	"encoding/json"

	etcd "go.etcd.io/etcd/client/v3"
)

type Service struct {
	Kind string `json:"Kind"`
	Addr string `json:"addr"`
}

type Query[T any] struct {
	*etcd.Client
}

func WithQuery[T any](cli *etcd.Client) *Query[T] {
	return &Query[T]{Client: cli}
}

func (x Query[T]) Query(key string) (docments []T, err error) {
	response, err := x.Get(context.Background(), key, etcd.WithPrefix())
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < len(response.Kvs); i++ {
		var doc T
		if err := json.Unmarshal(response.Kvs[i].Value, &doc); err != nil {
			return nil, err
		}
		docments = append(docments, doc)
	}
	return docments, nil
}
