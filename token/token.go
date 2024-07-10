package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Id          int64  `json:"id"`
	Version     string `json:"version"`
	Passport    string `json:"passport"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"e-mail"`
	jwt.RegisteredClaims
}

func GenToken(id int64, passport, phoneNumber, email string, key []byte) (token string, err error) {
	jwtToken := &Token{
		Id:          id,
		Version:     "1.0",
		Passport:    passport,
		PhoneNumber: phoneNumber,
		Email:       email,
	}
	return jwtToken.Marshal(key)
}

func (x *Token) Marshal(key []byte) (token string, err error) {
	x.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(8, 0, 0)),  // 过期时间3年
		IssuedAt:  jwt.NewNumericDate(time.Now().AddDate(-1, 0, 0)), // 签发时间
		NotBefore: jwt.NewNumericDate(time.Now().AddDate(-1, 0, 0)), // 生效时间
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, x)
	return jwtToken.SignedString(key)
}

func (x *Token) Unmarshal(token string, key []byte) (err error) {
	if len(token) <= 0 {
		return errors.New("len(token) <= 0")
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return err
	}
	v, ok := jwtToken.Claims.(*Token)
	if !ok {
		return errors.New("!ok = jwtToken.Claims.(*LoginToken)")
	}
	*x = *v
	return nil
}

func ParseToken(token string, key []byte) (*Token, error) {
	var t Token
	if err := t.Unmarshal(token, key); err != nil {
		return nil, err
	}
	return &t, nil
}

func GenerateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		//NextProtos:   []string{"http/1.1"},
		NextProtos: []string{"quic-echo-example"},
		MaxVersion: tls.VersionTLS13,
	}
}
