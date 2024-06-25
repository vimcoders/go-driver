package handle

import (
	"errors"
	"fmt"
	"go-driver/pb"

	"google.golang.org/protobuf/proto"
)

type Message []proto.Message

// 将一个来自底层的二进制流反序列化成一个对象
func (x Message) Unmarshal(req []byte) (proto.Message, error) {
	if len(req) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	var request Request = req
	kind := request.Kind()
	if kind >= uint16(len(x)) {
		return nil, fmt.Errorf("message id %v not registered", kind)
	}
	message := x[kind].ProtoReflect().New().Interface()
	if err := proto.Unmarshal(request.Message(), message); err != nil {
		return nil, err
	}
	return message, nil
}

// 将一个对象序列化成一个二进制流
func (x Message) Marshal(response proto.Message) ([]byte, error) {
	for i := uint16(0); i < uint16(len(x)); i++ {
		if proto.MessageName(response) != proto.MessageName(x[i]) {
			continue
		}
		return encode(i, response)
	}
	return nil, fmt.Errorf("message %s not registered", proto.MessageName(response))
}

// 定义所有的协议号
var Messages = Message{
	&pb.PingRequest{},
	&pb.PingResponse{},
	&pb.LoginRequest{},
	&pb.LoginResponse{},
}

// func (x *Protobuf) Unmarshal(b []byte) (proto.Message, proto.Message, error) {
// 	if len(b) < 2 {
// 		return nil, nil, errors.New("protobuf data too short")
// 	}
// 	id := binary.BigEndian.Uint16(b)
// 	if id >= uint16(len(x.Messages)) {
// 		return nil, nil, fmt.Errorf("message id %v not registered", id)
// 	}
// 	if id+1 >= uint16(len(x.Messages)) {
// 		return nil, nil, fmt.Errorf("message id %v not registered", id+1)
// 	}
// 	message := x.Messages[id].ProtoReflect().New().Interface()
// 	if err := proto.Unmarshal(b[2:], message); err != nil {
// 		return nil, nil, err
// 	}
// 	return message, x.Messages[id+1].ProtoReflect().New().Interface(), nil
// }

// func (x *Protobuf) Marshal(msg proto.Message) ([]byte, error) {
// 	for i := 0; i < len(x.Messages); i++ {
// 		if proto.MessageName(msg) != proto.MessageName(x.Messages[i]) {
// 			continue
// 		}
// 		data, err := proto.Marshal(msg)
// 		if err != nil {
// 			return nil, err
// 		}
// 		var header [2]byte
// 		binary.BigEndian.PutUint16(header[:], uint16(i))
// 		return append(header[:], data...), nil
// 	}
// 	return nil, fmt.Errorf("message %s not registered", proto.MessageName(msg))
// }

// func (x *Protobuf) Encode(origData []byte) ([]byte, error) {
// 	if len(x.Key) <= 0 {
// 		return origData, nil
// 	}
// 	block, err := aes.NewCipher(x.Key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	blockSize := block.BlockSize()
// 	origData = PKCS7Padding(origData, blockSize)
// 	blockMode := cipher.NewCBCEncrypter(block, x.Key[:blockSize])
// 	crypted := make([]byte, len(origData))
// 	blockMode.CryptBlocks(crypted, origData)
// 	return crypted, nil
// }

// func (x *Protobuf) Decode(crypted []byte) ([]byte, error) {
// 	if len(x.Key) <= 0 {
// 		return crypted, nil
// 	}
// 	block, err := aes.NewCipher(x.Key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	blockSize := block.BlockSize()
// 	blockMode := cipher.NewCBCDecrypter(block, x.Key[:blockSize])
// 	origData := make([]byte, len(crypted))
// 	blockMode.CryptBlocks(origData, crypted)
// 	origData = PKCS7UnPadding(origData)
// 	return origData, nil
// }

// func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
// 	padding := blockSize - len(ciphertext)%blockSize
// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(ciphertext, padtext...)
// }

// func PKCS7UnPadding(origData []byte) []byte {
// 	length := len(origData)
// 	unpadding := int(origData[length-1])
// 	return origData[:(length - unpadding)]
// }

// // ZeroUnPadding 去除填充字节的函数
// func ZeroUnPadding(data []byte) []byte {
// 	return bytes.TrimRightFunc(data, func(r rune) bool {
// 		return r == 0
// 	})
// }

// // ZeroPadding 填充字节的函数
// func ZeroPadding(data []byte, blockSize int) []byte {
// 	padding := blockSize - len(data)%blockSize
// 	fmt.Println("要填充的字节：", padding)
// 	slice1 := []byte{0}
// 	slice2 := bytes.Repeat(slice1, padding)
// 	return append(data, slice2...)
// }

// func NewProtobuf(message ...proto.Message) *Protobuf {
// 	return &Protobuf{Messages: append([]proto.Message{}, message...)}
// }
