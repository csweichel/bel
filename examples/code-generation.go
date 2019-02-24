package main

import (
	"os"

	"github.com/32leaves/bel"
)

type HelloWorld struct {
	Name    string
	Message string
}

func CodeGenCustomization() {
	extract, err := bel.Extract(HelloWorld{})
	if err != nil {
		panic(err)
	}

	f, err := os.Create("helloworld.ts")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// wrap in a namespace, add to the preamble and write to file
	err = bel.Render(extract,
		bel.GenerateAdditionalPreamble("// Hello World\n"),
		bel.GenerateNamespace("helloworld"),
		bel.GenerateOutputTo(f),
	)

	if err != nil {
		panic(err)
	}
}
