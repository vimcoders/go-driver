package handler

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-driver/app/proxy/driver"
	"go-driver/log"
	"go-driver/token"
)

func (x *Handler) PassportLogin(w driver.Response, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var request driver.PassportLoginRequest
	if err := json.Unmarshal(b, &request); err != nil {
		log.Error(err.Error())
		return
	}
	// TODO :: 校验用户名和密码是否符合规则
	hash := sha256.Sum256([]byte(request.Pwd))
	account := &driver.Account{
		Passport: request.Passport,
		Pwd:      fmt.Sprintf("%x", md5.Sum(hash[:])),
		Created:  time.Now(),
	}
	if err := x.Insert(&account); err != nil {
		log.Error(err.Error())
		return
	}
	jwtToken, err := token.GenToken(account.UserId, request.Passport, "", "", []byte(x.Token.Key))
	if err != nil {
		log.Error(err.Error())
		return
	}
	w.Write(driver.PassportLoginResponse{Token: jwtToken, Methods: x.Methods, Address: "127.0.0.1", Port: 9600})
}
