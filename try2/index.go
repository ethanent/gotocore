package main

import (
	"fmt"

	"github.com/ethanent/gotocore"
)

type user struct {
	Name string `g:"0,string"`
	Age  int    `g:"0,uint,8"`
}

func main() {
	ethan := &user{
		Name: "Ethan",
		Age:  18,
	}

	d := gotocore.Marshal(ethan)

	fmt.Println(d)
}
