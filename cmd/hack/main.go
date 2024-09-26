package main

import (
	"fmt"
	"os"

	"github.com/KScaesar/go-layout/pkg/utility"
)

func main() {
	var key string
	if len(os.Args) > 1 {
		key = os.Args[1]
	}
	hack := utility.Hack(key)
	fmt.Println(hack.Value())
}
