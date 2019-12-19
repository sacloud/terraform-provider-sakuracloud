package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func main() {
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		fmt.Printf("%s: %d\n", arg, schema.HashString(arg))
	}
}
