package main

import (
	"fmt"
	"os"

	"github.com/ombima56/go-ls-commands/listfiles"
)

func main() {
	args := os.Args[1:]
	path := "."

	if len(args) == 0 {
		listfiles.ListFiles(".", false, false)
		return
	}

	longFlag, allFlag, err := listfiles.ValidateFlags(args)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	listfiles.ListFiles(path, longFlag, allFlag)
}
