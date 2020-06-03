package main

import (
	"fmt"

	"github.com/ethanent/gotocore"
)

type user struct {
	Name string `g:"string"`
	Age  int    `g:"uint"`
}

func main() {
	ethan := &user{
		Name: "Ethan",
		Age:  18,
	}

	d := gotocore.Build(ethan)

	fmt.Println(d)
}
