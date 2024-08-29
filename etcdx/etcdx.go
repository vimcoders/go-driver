package etcdx

type Service struct {
	Kind      string `json:"kind"`
	Internet  string `json:"internet"`
	LocalAddr string `json:"localaddr"`
	Network   string `json:"network"`
}
