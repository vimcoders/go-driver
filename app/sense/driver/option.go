package driver

import (
	"flag"
	"go-driver/driver"
	"os"

	"gopkg.in/yaml.v3"
)

type Option = driver.Option

func ParseOption() (opt Option) {
	var fileName string
	flag.StringVar(&fileName, "option", "scene.conf", "scene.conf")
	flag.Parse()
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	if err := yaml.Unmarshal(ymalBytes, &opt); err != nil {
		panic(err.Error())
	}
	return opt
}
