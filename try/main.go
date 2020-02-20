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

	d, _, err := sch.Parse([]byte{1, 56, 2, 2, 44, 88, 7})

	fmt.Println(d)
	fmt.Println(err)
}
