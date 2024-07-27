package driver

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

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

type Mysql struct {
	Host string `ymal:"host"`
}

type Mongo struct {
	Host string `ymal:"host"`
	DB   string `ymal:"db"`
}

type Option struct {
	TCP      Addr     `yaml:"tcp"`
	HTTP     Addr     `yaml:"http"`
	QUIC     Addr     `yaml:"quic"`
	Mysql    Mysql    `yaml:"mysql"`
	Mongo    Mongo    `yaml:"mongo"`
	Etcd     Etcd     `yaml:"etcd"`
	Dingding Dingding `yaml:"dingding"`
	Telegram Telegram `yaml:"telegram"`
	Token    Token    `yaml:"token"`
}

func ReadOption() *Option {
	var fileName string
	flag.StringVar(&fileName, "conf", "./logic.conf", "logic.conf")
	flag.Parse()
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	var opt Option
	if err := yaml.Unmarshal(ymalBytes, &opt); err != nil {
		panic(err.Error())
	}
	return &opt
}
