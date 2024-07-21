// handle 按理说是应该在进程内自己实现，但是我们为了提高复用从外部引入
package driver

type Token struct {
	Key string `ymal:"key"`
}

type Addr struct {
	Internet string `ymal:"internet"`
	Port     string `ymal:"port"`
}

type Telegram struct {
	Token  string `ymal:"token"`
	ChatId string `ymal:"chat_id"`
}

type Dingding struct {
	Token  string `ymal:"token"`
	Secret string `ymal:"secret"`
}

type Etcd struct {
	Endpoints string `ymal:"endpoints"`
	UserName  string `ymal:"user_name"`
	Passwd    string `ymal:"passwd"`
	Version   string `ymal:"version"`
}

type YAML struct {
	TCP      Addr     `yaml:"tcp"`
	HTTP     Addr     `yaml:"http"`
	QUIC     Addr     `yaml:"quic"`
	Etcd     Etcd     `yaml:"etcd"`
	Dingding Dingding `yaml:"dingding"`
	Telegram Telegram `yaml:"telegram"`
	Token    Token    `yaml:"token"`
}
