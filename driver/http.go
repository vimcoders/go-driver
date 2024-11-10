// 不允许调用标准库外的包，防止循环引用
package driver

type PassportLoginRequest struct {
	Passport string `json:"Passport"`
	Pwd      string `json:"Pwd"`
}

type Metod struct {
	Id           int
	MethodName   string `json:"MethodName"`
	RequestName  string `json:"RequestName"`
	ResponseName string `json:"ResponseName"`
}

type PassportLoginResponse struct {
	Token   string  `json:"Token"`
	Methods []Metod `json:"Methods"`
	Address string  `json:"Address"`
	Port    int     `json:"Port"`
}
