package handler

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"time"

	"go-driver/app/proxy/driver"
	"go-driver/token"
)

func (x *Handler) PassportLogin(ctx context.Context, req *driver.PassportLoginRequest) (*driver.PassportLoginResponse, error) {
	fmt.Println(ctx.Value("Authorization"))
	// TODO :: 校验用户名和密码是否符合规则
	hash := sha256.Sum256([]byte(req.Pwd))
	account := &driver.Account{
		Passport: req.Passport,
		Pwd:      fmt.Sprintf("%x", md5.Sum(hash[:])),
		Created:  time.Now(),
	}
	if err := x.Insert(&account); err != nil {
		return &driver.PassportLoginResponse{}, err
	}
	jwtToken, err := token.GenToken(account.UserId, req.Passport, "", "", []byte(x.Token.Key))
	if err != nil {
		return &driver.PassportLoginResponse{}, err
	}
	return &driver.PassportLoginResponse{Token: jwtToken, Address: "127.0.0.1", Port: 9600}, nil
}
