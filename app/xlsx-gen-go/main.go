package main

import "go-driver/app/xlsx-gen-go/generator"

func main() {
	g, err := generator.NewGenerator()
	if err != nil {
		panic(err)
	}
	if err := g.Gen(); err != nil {
		panic(err)
	}
}
