package main

import (
	"fmt"

	"github.com/ello/ello-go/common/util"
)

func main() {
	v, err := util.ValidateInt("123", 0)
	fmt.Printf("THIS IS STREAMZ with example values of %v %v", v, err)
}
