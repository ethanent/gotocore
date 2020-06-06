package main

import (
	"crypto/rand"
	"fmt"

	"github.com/ethanent/gotocore"
)

type address struct {
	Number      int    `g:"0,varint"`
	Street      string `g:"1,string"`
	City        string `g:"2,string"`
	AdminRegion string `g:"3,string"`
}

type user struct {
	Name string   `g:"0,string"`
	Addr *address `g:"1"`
	Age  int      `g:"2,uint,16"`
}

func main() {
	ethan := &user{
		Name: "Ethan",
		Addr: &address{
			Number:      24,
			Street:      "Bayview Avenue",
			City:        "San Francisco",
			AdminRegion: "California",
		},
		Age: 2552,
	}

	d, err := gotocore.Marshal(ethan)

	if err != nil {
		panic(err)
	}

	fmt.Println("Prenoise length:", len(d))

	// Add some fun noise to check that the right amount of bytes is read

	for i := 0; i < 100; i++ {
		r := make([]byte, 1, 1)
		_, err := rand.Reader.Read(r)

		if err != nil {
			panic(err)
		}

		d = append(d, r[0])
	}

	// Print

	fmt.Println(d)

	ex := &user{}

	c, err := gotocore.Unmarshal(d, ex)

	if err != nil {
		panic(err)
	}

	fmt.Println("Read", c, "/", len(d))

	fmt.Println(ex)
}
