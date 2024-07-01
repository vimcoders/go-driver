package grpcx

import (
	"bufio"
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/proto"
)

type Message []byte

func decode(b *bufio.Reader) (Message, error) {
	headerBytes, err := b.Peek(4)
	if err != nil {
		return nil, err
	}
	header := binary.BigEndian.Uint32(headerBytes)
	if int(header) > b.Size() {
		return nil, fmt.Errorf("header %v too long", header)
	}
	iMessage, err := b.Peek(int(header))
	if err != nil {
		return nil, err
	}
	if _, err := b.Discard(len(iMessage)); err != nil {
		return nil, err
	}
	return iMessage, nil
}

func encode(seq, ack uint32, kind uint16, iMessage proto.Message) (Message, error) {
	b, err := proto.Marshal(iMessage)
	if err != nil {
		return nil, err
	}
	var message [14]byte
	binary.BigEndian.PutUint32(message[:], uint32(len(message)+len(b)))
	binary.BigEndian.PutUint32(message[4:], seq)
	binary.BigEndian.PutUint32(message[8:], ack)
	binary.BigEndian.PutUint16(message[12:], uint16(kind))
	return append(message[:], b...), nil
}

func (x Message) seq() uint32 {
	return binary.BigEndian.Uint32(x[4:])
}

func (x Message) ack() uint32 {
	return binary.BigEndian.Uint32(x[8:])
}

func (x Message) method() uint16 {
	return binary.BigEndian.Uint16(x[12:])
}

func (x Message) payload() []byte {
	return x[14:]
}

func (x Message) clone() Message {
	clone := make(Message, len(x))
	copy(clone, x)
	return clone
}
