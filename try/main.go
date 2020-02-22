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

	sch.Components = append(sch.Components, protocore.Component{
		Name: "tstr",
		Kind: protocore.String,
	})

	d, _, err := sch.Parse([]byte{128, 56, 2, 44, 88, 7, 0, 6, 56, 69, 69, 69, 42, 0, 0, 15, 72, 101, 89, 32, 84, 72, 69, 114, 101, 33, 32, 51, 53, 52, 54})

	fmt.Println(d)
	fmt.Println(err)

	b := sch.Build(map[string]interface{}{
		"uname": -56,
		"tsts":  481324,
		"tbuf":  []byte{56, 69, 69, 69, 42, 0},
		"tstr":  "HeY THEre! 3546",
	})

	fmt.Println(b)
}
