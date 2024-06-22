package router

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"go-driver/app/proxy/handler"
	"go-driver/driver"
	"go-driver/log"
)

type Router struct {
	trees map[string]func(w driver.Response, r *http.Request)
}

func NewRouter(handler *handler.Handler) http.Handler {
	return &Router{
		trees: map[string]func(w driver.Response, r *http.Request){
			"/api/v1/passport/login": handler.PassportLogin,
		},
	}
}

func (x *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(fmt.Sprintf("%s", e))
			debug.PrintStack()
		}
		r.Body.Close()
	}()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	log.Info(r.URL.Path)
	if v, ok := x.trees[r.URL.Path]; ok {
		v(&driver.ResponseWriter{W: w}, r)
	}
}
