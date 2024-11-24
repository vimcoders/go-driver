package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go-driver/log"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Context struct {
	context.Context
	value map[any]any
}

func WithContext(ctx context.Context) *Context {
	return &Context{Context: ctx, value: map[any]any{}}
}

func (x *Context) Value(k any) any {
	if v, ok := x.value[k]; ok {
		return v
	}
	return nil
}

func (x *Context) WithValue(k, v any) {
	x.value[k] = v
}

func (x *Handler) NewRouter() http.Handler {
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
	log.Info(r.URL.Path)
	paths := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")
	if len(paths) <= 0 {
		return
	}
	var methodName string
	for i := 0; i < len(paths); i++ {
		methodName += cases.Title(language.English).String(paths[i])
	}
	ctx := WithContext(context.Background())
	for k, v := range r.Header {
		if len(v) <= 0 {
			continue
		}
		ctx.WithValue(k, v[0])
	}
	for i := 0; i < len(x.ServiceDesc.Methods); i++ {
		if methodName != x.ServiceDesc.Methods[i].MethodName {
			continue
		}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		dec := func(in any) error {
			if err := json.Unmarshal(b, in); err != nil {
				return err
			}
			return nil
		}
		reply, err := x.ServiceDesc.Methods[i].Handler(x, ctx, dec, nil)
		if err != nil {
			log.Error(err.Error())
		}
		response, err := json.Marshal(reply)
		if err != nil {
			log.Error(err.Error())
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Expose-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if _, err := w.Write(response); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
