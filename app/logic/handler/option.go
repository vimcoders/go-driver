package handler

import (
	"flag"
	"go-driver/app/logic/driver"
	"os"

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
