package handler

import (
	"context"
	"go-driver/app/proxy/driver"
	"go-driver/pb"
	"go-driver/sqlx"
	"net/http"

	"google.golang.org/protobuf/proto"
)

var handler *Handler

type Handler struct {
	Option
	sqlx.Client
	trees   map[string]func(w driver.Response, r *http.Request)
	Methods []driver.Metod
}

func MakeHandler(ctx context.Context) *Handler {
	var methods []driver.Metod = make([]driver.Metod, len(pb.Parkour_ServiceDesc.Methods))
	for i := 0; i < len(pb.Parkour_ServiceDesc.Methods); i++ {
		method := pb.Parkour_ServiceDesc.Methods[i]
		dec := func(in any) error {
			methods[i].RequestName = string(proto.MessageName(in.(proto.Message)).Name())
			return nil
		}
		resp, _ := method.Handler(&pb.UnimplementedParkourServer{}, context.Background(), dec, nil)
		methods[i].Id = i
		methods[i].MethodName = method.MethodName
		methods[i].ResponseName = string(proto.MessageName(resp.(proto.Message)).Name())
	}
	h := &Handler{Methods: methods}
	if err := h.Parse(); err != nil {
		panic(err)
	}
	if err := h.Connect(ctx); err != nil {
		panic(err)
	}
	handler = h
	return handler
}
