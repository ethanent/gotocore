package main

import (
	"fmt"

	"github.com/ethanent/protocore-go"
)

func main() {
	sch := protocore.Schema{}

	sch.Components = append(sch.Components, protocore.Component{
		Name: "uname",
		Kind: protocore.Varint,
	})

	sch.Components = append(sch.Components, protocore.Component{
		Name: "tsts",
		Kind: protocore.Varint,
	})

	d, _, err := sch.Parse([]byte{128, 56, 2, 44, 88, 7})

	fmt.Println(d)
	fmt.Println(err)

	b := sch.Build(map[string]interface{}{
		"uname": -56,
		"tsts":  481324,
	})

	fmt.Println(b)
}
