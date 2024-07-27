package handler

import (
	"fmt"
	"go-driver/driver"
	"go-driver/log"
	"net/http"
	"runtime/debug"
)

func (x *Handler) NewRouter() http.Handler {
	x.trees = map[string]func(w driver.Response, r *http.Request){
		"/api/v1/passport/login": x.PassportLogin,
	}
	return x
}

func (x *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
