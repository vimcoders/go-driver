package driver

type PassportLoginRequest struct {
	Passport string `json:"passport"`
	Pwd      string `json:"pwd"`
}

type PassportLoginResponse struct {
	Token string `json:"token"`
}
