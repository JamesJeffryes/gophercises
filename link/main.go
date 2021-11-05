package main

import (
	"fmt"
	"github.com/jamesjeffryes/gophercises/link/link"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a path to a HTML file to parse")
		os.Exit(1)
	}
	fmt.Println(link.ParseFile(os.Args[1]))
}
