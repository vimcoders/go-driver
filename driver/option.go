package driver

import (
	"strings"
)

type Token struct {
	Key string `ymal:"key"`
}

type Endpoints struct {
	Wan string `ymal:"wan"`
	Lan string `ymal:"lan"`
}

func (x *Endpoints) LAN() string {
	return x.Lan
}

func (x *Endpoints) WAN() string {
	return x.Wan
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

func (x *Etcd) Join(prefix ...string) string {
	return x.Version + strings.Join(prefix, "/")
}

type Mysql struct {
	Host string `ymal:"host"`
}

type Mongo struct {
	Host string `ymal:"host"`
	DB   string `ymal:"db"`
}

type Option struct {
	TCP      Endpoints `yaml:"tcp"`
	HTTP     Endpoints `yaml:"http"`
	QUIC     Endpoints `yaml:"quic"`
	Etcd     Etcd      `yaml:"etcd"`
	Dingding Dingding  `yaml:"dingding"`
	Telegram Telegram  `yaml:"telegram"`
	Token    Token     `yaml:"token"`
	Mysql    Mysql     `yaml:"mysql"`
	Mongo    Mongo     `yaml:"mongo"`
}
