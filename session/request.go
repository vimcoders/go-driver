package session

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

type Request []byte

func parseRequest(b *bufio.Reader) (request Request, err error) {
	if err := request.parse(b); err != nil {
		return nil, err
	}
	return request, nil
}

func NewRequest(kind uint16) Request {
	var b [6]byte
	var header Request = b[:]
	binary.BigEndian.PutUint32(header[:], uint32(len(header)))
	binary.BigEndian.PutUint16(header[4:], uint16(kind))
	return header
}

func (x *Request) parse(b *bufio.Reader) error {
	headerBytes, err := b.Peek(4)
	if err != nil {
		return err
	}
	header := binary.BigEndian.Uint32(headerBytes)
	if int(header) > b.Size() {
		return fmt.Errorf("header %v too long", header)
	}
	request, err := b.Peek(int(header))
	if err != nil {
		return err
	}
	if _, err := b.Discard(len(request)); err != nil {
		return err
	}
	*x = request
	return nil
}

func (x Request) Kind() uint16 {
	return binary.BigEndian.Uint16(x[4:])
}

func (x Request) Message() []byte {
	return x[6:]
}

func (x Request) Reply() uint16 {
	return x.Kind() + 1
}
