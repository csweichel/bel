package main

import (
	"github.com/32leaves/bel"
)

// WithDocumentation demonstrates how to extract documentation
func WithDocumentation() {
	handler, err := bel.NewParsedSourceDocHandler(".", "github.com/32leaves")
	if err != nil {
		panic(err)
	}
	// add the main package in exmaples/
	if err := handler.AddToIndex("examples", ""); err != nil {
		panic(err)
	}

	ts, err := bel.Extract((*UserService)(nil), bel.WithDocumentation(handler), bel.FollowStructs)
	if err != nil {
		panic(err)
	}

	err = bel.Render(ts)
	if err != nil {
		panic(err)
	}
}

func init() {
	examples["with-documentation"] = WithDocumentation
}
