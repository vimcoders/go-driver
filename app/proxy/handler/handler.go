package handler

import (
	"context"
	"go-driver/app/proxy/driver"
	"go-driver/sqlx"
	"net/http"
)

var handler *Handler

type Handler struct {
	Option
	sqlx.Client
	trees map[string]func(w driver.Response, r *http.Request)
}

func MakeHandler(ctx context.Context) *Handler {
	h := &Handler{}
	if err := h.Parse(); err != nil {
		panic(err)
	}
	if err := h.Connect(ctx); err != nil {
		panic(err)
	}
	handler = h
	return handler
}
