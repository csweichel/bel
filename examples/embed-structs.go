package main

import (
	"github.com/32leaves/bel"
)

type ReferencedStruct struct {
	FieldA string
	FieldB int32
}

type HasAReference struct {
	SomeRef ReferencedStruct
}

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
