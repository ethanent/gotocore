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

	sch.Components = append(sch.Components, protocore.Component{
		Name: "tbuf",
		Kind: protocore.Buffer,
	})

	d, _, err := sch.Parse([]byte{128, 56, 2, 44, 88, 7, 0, 6, 56, 69, 69, 69, 42, 0})

	fmt.Println(d)
	fmt.Println(err)

	b := sch.Build(map[string]interface{}{
		"uname": -56,
		"tsts":  481324,
		"tbuf":  []byte{56, 69, 69, 69, 42, 0},
	})

	fmt.Println(b)
}
