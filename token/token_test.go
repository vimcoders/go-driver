package token_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"go-driver/token"
	"os"
	"testing"
)

func Encrypt(plantText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	plantText = PKCS7Padding(plantText, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, key)
	ciphertext := make([]byte, len(plantText))
	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes) //选择加密算法
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, keyBytes)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	plantText = PKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func TestPush(t *testing.T) {
	// var jwtToken jwt.FightToken
	// jwtToken.Unmarshal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6IjEiLCJVc2VySWQiOjUyMywiUm9sZUlkIjoyMDAxMDAxMDAwMDExLCJleHAiOjE5NTA3NDUzODMsIm5iZiI6MTY2Njc0ODU4MywiaWF0IjoxNjY2NzQ4NTgzfQ.7KWpb6e6cWtikSA7vp9pUB02OrW1qRTnOXg_xmjASdM")
	// t.Log(jwtToken.UserId, jwtToken.RoleId)
}

func TestPrivateKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
		return
	}
	file, err := os.Create("parkour.pem")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	err = pem.Encode(file, block)
	if err != nil {
		t.Error(err)
	}
}

func TestGenPrivateKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Error(err)
		return
	}
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	buffer := bytes.NewBuffer(nil)
	err = pem.Encode(buffer, block)
	if err != nil {
		t.Error(err)
	}
	publicBuffer := bytes.NewBuffer(nil)
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}
	err = pem.Encode(publicBuffer, publicBlock)
	if err != nil {
		t.Error(err)
	}
	t.Log(buffer.String(), publicBuffer.String())
}

func TestUnmarshalPrivateKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Error(err)
		return
	}
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	buffer := bytes.NewBuffer(nil)
	err = pem.Encode(buffer, block)
	if err != nil {
		t.Error(err)
	}
	decodeBlock, _ := pem.Decode(buffer.Bytes())
	key, err := x509.ParsePKCS1PrivateKey(decodeBlock.Bytes)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(buffer.String(), key)
}

func TestParseECPrivateKey(t *testing.T) {
	derKey := []byte(`-----BEGIN EC PRIVATE KEY-----
	MHUCAQEEIU0bf1T5F1MLSHsQ28bg1SIyN/pV1HLQJdR9SdoNULdNjKAHBgUrgQQACqFEA
	0IABFdjz/axFdsteRT4TwS3UDnBZSAeToecWN5u/ZkDDRZ+7uhFOeqmkQTAPPrdK+MQ4rHb
	giPXQPkKl9pFIwVqp0o=
	-----END EC PRIVATE KEY-----`)
	privKey, err := x509.ParseECPrivateKey(derKey)
	if err != nil {
		panic(err)
	}
	if err != nil {
		t.Error(err)
	}
	t.Log(privKey)
}

func TestParsePrivateKey(t *testing.T) {
	block, _ := pem.Decode([]byte(`-----BEGIN PRIVATE KEY-----
	MIIEpQIBAAKCAQEApW0xeOLzRck5apE3IP72nipk/FkyIEuEGm838H8hzEmPo/ip
	TcmfTzRXwgC4XCRfoVOt1XBKuta6Avi2Dfgh5GXDzw7WnwfzkO8OyKmTJsL/HcnN
	EuUTCtP7mZwUgosrCVy8UqjHc82gT9ZSBK7wILblNYYwL5BEK8ZEWrG7vJE0/70l
	BOXkutvflllWYrHNXCJHQSccFLYxNoExCKkOKTNUZpzKyDanPHcyaZM+kWtTMR6k
	bgmVGFXcig2d0RjdvnsI0CVTpSL8QViEI1OFXxd8CEkvocUeGytmzVfCXcPu1F4g
	+IqKnxRj0/uz8kMDOAQlvTP3iuvRAoD/OFx30wIDAQABAoIBAQCVbk2CJYAbSenT
	mdlytN2Rgjo2uVvOUGjEeDLPzAd7wfc+5yAIZFjD80RSutPOaAz6bdxZMVZP8CeX
	B5NsivgSmNqH759viH88LLXuDUAfg4VwIxpcNxE8dsCPwa3FPnFhw6NaB5wjv1tQ
	wwjTsjK3Wn8yGkTssiTiZfbY9jPf4NbcSTtB3kwtABwIIDxaPWSv1ps6aV5cUCMt
	QHWjWMDvinkqWX6zbdD0ULkotfER1oGXpKSS92hOzu9HE8GLQvd/UXUNszkaJsvG
	GfVvWnix0badNMVQp0c1GZLVOGYhf1VkLFx49cs3QG3H21njzP1uq395BV1tmSiH
	6cJ6EQYBAoGBAMJu02jrOZQ8g/tUucgAPDI/Mi8O1BwH07SEaGggQoqxvgjpf1/O
	UMAEmraWFxTUATK8BwQQBp9+KP7dkhJZklBg/iWpHmdVRzXV4Oav7IYXH2YtPehZ
	ViBT8hBGmr0PlUNmyiGNVnb+cjwRZSODZNfJXHr0Pyv1uNYIXT5RZ9bVAoGBANnP
	C2W0RQX/QkP6g07BIyMU7qmycVYY3GaRAjd0p2T8zjwJ50XVuigbw2ECVHmwX82B
	7zx+jEBMuKLn+ODZwTZ0JKb2JMiNb/Pj+FOENzjLG1XcPX+3xhEHSi+T+V+xYD+d
	uekd6MfgymbTJ3BVPPQZg309LFOG2TcqxQ7K6zgHAoGBAIXLkbs/Mu4o/oFy+i0A
	zGufRT9QqvFnCW3NN7N/j4q1aRnk4/vfk32vLW+7tMJmaTSqYwGOraAPRtKrUhtC
	fAbH19u+ludwrYIEXbEhGlfjjX3YYCOFZlj0qzw7+btj/8jT8QBJrFhSG/Xt2nUn
	s7syG2uYq+fqPXk7ZD6/8f7JAoGALk/u3XZCQu8uuOOYbfN1NC1sPdr6bFMm8gwd
	S4tbWIbEl1GHwnqadZLJrWPgcGuHQ1xAcT17NuTZUZI/ghfrFFgHvxSRZ69jQZmU
	oLV5RHMzYcNNtE1wKQjCxnERUj6V95DjCeVZLL7oaoq1VRZaupB+O+/47925bBiF
	BAszjpsCgYEAj00CMS/JSFzf8hBgL2dZyB6bQRKyz49vUpsqX/jKIFbr6o/s7+Ot
	HWm58MiXwr67BYi4+WFuhDiNdiz4cuhxqFnIQ/CpZxellamjrVvVeHWYAn3qp1gz
	CvcdmLIt3b1D+oBf6rt485ku1XGYIe+zkjLEYt/1tkhYAlqTD++x+gM=
	-----END PRIVATE KEY-----`))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Error(err)
	}
	t.Log(key)
}

// func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
// 	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
// 	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

// 	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
// 	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

// 	return string(pemEncoded), string(pemEncodedPub)
// }

// func decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
// 	block, _ := pem.Decode([]byte(pemEncoded))
// 	x509Encoded := block.Bytes
// 	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

// 	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
// 	x509EncodedPub := blockPub.Bytes
// 	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
// 	publicKey := genericPublicKey.(*ecdsa.PublicKey)

// 	return privateKey, publicKey
// }

// func test() {
// 	privateKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
// 	publicKey := &privateKey.PublicKey

// 	encPriv, encPub := encode(privateKey, publicKey)

// 	fmt.Println(encPriv)
// 	fmt.Println(encPub)

// 	priv2, pub2 := decode(encPriv, encPub)

// 	if !reflect.DeepEqual(privateKey, priv2) {
// 		fmt.Println("Private keys do not match.")
// 	}
// 	if !reflect.DeepEqual(publicKey, pub2) {
// 		fmt.Println("Public keys do not match.")
// 	}
// }

func TestJWT(t *testing.T) {
	token := &token.Token{
		Id:          1,
		Version:     "test",
		Passport:    "passportxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		Email:       "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		PhoneNumber: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	}
	jwtToken, _ := token.Marshal([]byte("123"))
	t.Log(jwtToken)
}
