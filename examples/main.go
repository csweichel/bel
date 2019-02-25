package main

import (
	"fmt"
	"os"
)

var examples = make(map[string]func())

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	name := os.Args[1]
	example, ok := examples[name]
	if !ok {
		fmt.Printf("unknown example: %s\n", name)
		return
	}

	example()
}

func printUsage() {
	fmt.Printf("usage: %s <example>\nAvailable examples:\n", os.Args[0])
	for k := range examples {
		fmt.Printf("  %s\n", k)
	}
}
