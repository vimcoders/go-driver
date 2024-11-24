package handler

import (
	"context"
	"go-driver/pb"
	"go-driver/sqlx"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	Option
	sqlx.Client
	Methods []*pb.Method
	grpc.ServiceDesc
	pb.ProxyServer
}

func MakeHandler(ctx context.Context) *Handler {
	var methods []*pb.Method = make([]*pb.Method, len(pb.Parkour_ServiceDesc.Methods))
	for i := 0; i < len(pb.Parkour_ServiceDesc.Methods); i++ {
		method := pb.Parkour_ServiceDesc.Methods[i]
		var newMethod pb.Method
		dec := func(in any) error {
			newMethod.RequestName = string(proto.MessageName(in.(proto.Message)).Name())
			return nil
		}
		resp, _ := method.Handler(&pb.UnimplementedParkourServer{}, context.Background(), dec, nil)
		newMethod.Id = int32(i)
		newMethod.MethodName = method.MethodName
		newMethod.ResponseName = string(proto.MessageName(resp.(proto.Message)).Name())
		methods[i] = &newMethod
	}
	h := &Handler{Methods: methods, ServiceDesc: pb.Proxy_ServiceDesc}
	if err := h.Parse(); err != nil {
		panic(err)
	}
	if err := h.Connect(ctx); err != nil {
		panic(err)
	}
	return h
}
