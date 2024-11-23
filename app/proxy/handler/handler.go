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

type Method struct {
	MethodName string
	method     reflect.Value
	req        interface{}
}

func (x *Method) NewRequest() interface{} {
	t := reflect.TypeOf(x.req).Elem()
	return reflect.New(t).Interface()
}

func (x *Handler) Call(ctx context.Context, methodName string, dec func(req interface{}) error) (interface{}, error) {
	method := reflect.ValueOf(x).MethodByName(methodName)
	if ok := method.IsValid(); !ok {
		return nil, errors.New(fmt.Sprintf("method.IsValid(); !ok  %s", methodName))
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
	if err := result[1].Interface().(error); err != nil {
		return nil, err
	}
	response := result[0].Interface()
	if response == nil {
		return nil, errors.New("response == nil")
	}
	return response, nil
}

var handler *Handler

type Handler struct {
	Option
	sqlx.Client
	//Methods []driver.Metod
	Metohods []Method
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
	h := &Handler{}
	if err := h.Parse(); err != nil {
		panic(err)
	}
	if err := h.Connect(ctx); err != nil {
		panic(err)
	}
	t, v := reflect.TypeOf(h), reflect.ValueOf(h)
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		mt := method.Type()
		if mt.NumIn() != 2 {
			continue
		}
		if mt.NumOut() != 2 {
			continue
		}
		req := reflect.New(mt.In(1).Elem()).Interface()
		h.Metohods = append(h.Metohods, Method{
			MethodName: t.Method(i).Name,
			method:     method,
			req:        req,
		})
	}
	handler = h
	return handler
}
