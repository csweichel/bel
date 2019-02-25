package main

import (
	"github.com/32leaves/bel"
)

// ReferencedStruct is a struct referenced by HasAReference
type ReferencedStruct struct {
	FieldA string
	FieldB int32
}

// HasAReference is a struct with a reference
type HasAReference struct {
	SomeRef ReferencedStruct
}

// EmbedStructs demonstrates the EmbedStructs config option
func EmbedStructs() {
	ts, err := bel.Extract(HasAReference{}, bel.EmbedStructs)
	if err != nil {
		panic(err)
	}

	err = bel.Render(ts)
	if err != nil {
		panic(err)
	}
}

func init() {
	examples["embed-structs"] = EmbedStructs
}
