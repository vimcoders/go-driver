package handler

import "google.golang.org/protobuf/proto"

type Scene interface {
	Kind() uint16
	Join(Session) error
	Leave(proto.Message) error
	Broadcast(proto.Message) error
	Push(uint64, proto.Message) error
	PlayerCount() uint16
}
