package main

import (
	"fmt"

	"github.com/wedojava/enit"
)

func main() {
	v := enit.File("my-fake-key", ".secrets")
	err := v.Set("demo_key", "some crazy value")
	if err != nil {
		panic(err)
	}
	plain, err := v.Get("demo_key")
	if err != nil {
		panic(err)
	}
	fmt.Println("Plain:", plain)
}
