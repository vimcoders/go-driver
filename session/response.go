package session

import (
	"encoding/binary"

	"google.golang.org/protobuf/proto"
)

type Response []byte

func NewResponse(kind uint16, message proto.Message) (Response, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	var header [6]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(header)+len(b)))
	binary.BigEndian.PutUint16(header[4:], uint16(kind))
	return append(header[:], b...), nil
}
