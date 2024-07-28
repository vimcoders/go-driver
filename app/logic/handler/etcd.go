package handler

import (
	"context"
	"go-driver/log"
)

func (x *Handler) Watch(ctx context.Context) {
	for ev := range x.Client.Watch(ctx, x.Etcd.Join()) {
		log.Info(ev.Events)
	}
}
