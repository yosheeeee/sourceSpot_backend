package main

import (
	"fmt"

	"github.com/yosheeeee/sourceSpot_baackend/initialize"
)

func main() {
	if err := initialize.InitializeApp("../../config.yaml"); err != nil {
		fmt.Print(err)
	}
}
