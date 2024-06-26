// 不允许调用标准库外的包，防止循环引用
package driver

type PassportLoginRequest struct {
	Passport string `json:"passport"`
	Pwd      string `json:"pwd"`
}

type PassportLoginResponse struct {
	Token string `json:"token"`
}
