package etcdx

type Service struct {
	Kind string `json:"kind"`
	WAN  string `json:"wan"`
	LAN  string `json:"lan"`
}
