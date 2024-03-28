package message

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/vimcoders/go-driver/pb"
	"google.golang.org/protobuf/proto"
)

var Messages = []proto.Message{
	&pb.LoginRequest{},
	&pb.LoginResponse{},
}

type Protobuf struct {
	Messages []proto.Message
	Key      []byte
}

func (x *Protobuf) Register(message proto.Message) {
	for i := 0; i < len(x.Messages); i++ {
		if proto.MessageName(x.Messages[i]) == proto.MessageName(message) {
			panic(fmt.Sprintf("msg %s is already registered", proto.MessageName(message)))
		}
	}
	x.Messages = append(x.Messages, message)
}

func (x *Protobuf) Unmarshal(b []byte) (proto.Message, error) {
	if len(b) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	id := binary.BigEndian.Uint16(b)
	if id >= uint16(len(x.Messages)) {
		return nil, fmt.Errorf("message id %v not registered", id)
	}
	decode, err := x.Decode(b[2:])
	if err != nil {
		return nil, err
	}
	message := x.Messages[id].ProtoReflect().New().Interface()
	if err := proto.Unmarshal(decode, message); err != nil {
		return nil, err
	}
	return message, nil
}

func (x *Protobuf) Marshal(msg proto.Message) ([]byte, error) {
	for i := 0; i < len(x.Messages); i++ {
		if proto.MessageName(msg) != proto.MessageName(x.Messages[i]) {
			continue
		}
		data, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}
		encode, err := x.Decode(data)
		if err != nil {
			return nil, err
		}
		var header [2]byte
		binary.BigEndian.PutUint16(header[:], uint16(i))
		return append(header[:], encode...), nil
	}
	return nil, fmt.Errorf("message %s not registered", proto.MessageName(msg))
}

func (x *Protobuf) Encode(origData []byte) ([]byte, error) {
	if len(x.Key) <= 0 {
		return origData, nil
	}
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, x.Key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func (x *Protobuf) Decode(crypted []byte) ([]byte, error) {
	if len(x.Key) <= 0 {
		return crypted, nil
	}
	block, err := aes.NewCipher(x.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, x.Key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// ZeroUnPadding 去除填充字节的函数
func ZeroUnPadding(data []byte) []byte {
	return bytes.TrimRightFunc(data, func(r rune) bool {
		return r == 0
	})
}

// ZeroPadding 填充字节的函数
func ZeroPadding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	fmt.Println("要填充的字节：", padding)
	slice1 := []byte{0}
	slice2 := bytes.Repeat(slice1, padding)
	return append(data, slice2...)
}

func NewProtobuf(messages ...proto.Message) *Protobuf {
	return &Protobuf{Messages: append([]proto.Message{}, messages...)}
}
