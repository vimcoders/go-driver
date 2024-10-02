// 不允许调用标准库外的包，防止循环引用
package driver

type PassportLoginRequest struct {
	Passport string `json:"passport"`
	Pwd      string `json:"pwd"`
}

type Metod struct {
	Id       int
	Name     string
	Request  string
	Response string
}

type PassportLoginResponse struct {
	Token   string  `json:"token"`
	Methods []Metod `json:"methods"`
}
