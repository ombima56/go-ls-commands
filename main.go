package main

import (
	"fmt"
	"os"

	"github.com/ombima56/go-ls-commands/listfiles"
)

func main() {
	dir := "."
	showDetails := false
	includeHidden := false
	reverseOrder := false
	sortByTime := false
	recursive := false

	for i := 1; i < len(os.Args); i++ {
		args := os.Args[i]
		switch args {
		case "-l":
			showDetails = true
		case "-R":
			recursive = true
		case "-a":
			includeHidden = true
		case "-r":
			reverseOrder = true
		case "-t":
			sortByTime = true
		default:
			dir = args
		}
	}

	if err := listfiles.ListFiles(dir, showDetails, recursive, includeHidden, reverseOrder, sortByTime); err != nil {
		fmt.Println("Error:", err)
	}
}