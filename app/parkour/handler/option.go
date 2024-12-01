package handler

import (
	"flag"
	"os"

	"github.com/vimcoders/go-driver/app/parkour/driver"

	"gopkg.in/yaml.v3"
)

type Option = driver.Option

func (x *Handler) Parse() error {
	var fileName string
	flag.StringVar(&fileName, "option", "logic.conf", "logic.conf")
	flag.Parse()
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(ymalBytes, &x.Option); err != nil {
		return err
	}
	return nil
}
