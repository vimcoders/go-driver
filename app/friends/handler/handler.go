package handler

import (
	etcd "go.etcd.io/etcd/client/v3"
)

type Handler struct {
	*etcd.Client
}
