package main

import (
	"fmt"
	"os"

	"github.com/KScaesar/go-layout/pkg/utility"
)

var key string

func main() {
	if len(key) == 0 && len(os.Args) > 1 {
		key = os.Args[1]
	}
	hack := utility.Hack(key)
	fmt.Println(hack.Value())
}
