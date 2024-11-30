package handler

import (
	"context"
	"fmt"
	"go-driver/log"
	"net/http"
	"runtime/debug"
)

func (x *Handler) ListenAndServe(ctx context.Context) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
	}()
	srv := &http.Server{
		Addr:    x.HTTP.Internet,
		Handler: x.NewRouter(),
	}
	defer srv.Close()
	if err := srv.ListenAndServe(); err != nil {
		log.Errorf("listen: %s", err.Error())
	}
}
