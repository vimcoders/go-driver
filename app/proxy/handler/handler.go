package handler

import (
	"context"
	"errors"
	"fmt"
	"go-driver/app/proxy/driver"
	"go-driver/pb"
	"go-driver/sqlx"
	"reflect"

	"google.golang.org/protobuf/proto"
)

func (x *Handler) Call(ctx context.Context, methodName string, dec func(req interface{}) error) (interface{}, error) {
	method := reflect.ValueOf(x).MethodByName(methodName)
	if ok := method.IsValid(); !ok {
		return nil, fmt.Errorf("method.IsValid(); !ok  %s", methodName)
	}
	mt := method.Type()
	if mt.NumIn() != 2 {
		return nil, errors.New("mt.NumIn() != 2")
	}
	if mt.NumOut() != 2 {
		return nil, errors.New("mt.NumOut() != 2")
	}
	req := reflect.New(mt.In(1).Elem()).Interface()
	if err := dec(req); err != nil {
		return nil, err
	}
	result := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
	if len(result) != 2 {
		return nil, errors.New("len(result) != 2")
	}
	response := result[0].Interface()
	if response == nil {
		return nil, errors.New("response == nil")
	}
	if err := result[1].Interface(); err != nil {
		return nil, err.(error)
	}
	return response, nil
}

var handler *Handler

type Handler struct {
	Option
	sqlx.Client
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
