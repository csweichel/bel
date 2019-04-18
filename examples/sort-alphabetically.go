package main

import (
	"github.com/32leaves/bel"
)

// StructB is a structure with non-alphabetically sorted fields
type StructB struct {
	FieldC string
	FieldA string
	FieldB string
}

// StructA is a structure with non-alphabetically sorted fields
type StructA struct {
	FieldB string
	FieldC int
	FieldA string
}

// StructC uses StructA and StructB
type StructC struct {
	BMember StructB
	AMember StructA
}

// SortAlphabetically demonstrates the use of the bel.SortAlphabetically config option
func SortAlphabetically() {
	ts, err := bel.Extract(StructC{}, bel.SortAlphabetically, bel.FollowStructs)
	if err != nil {
		panic(err)
	}

	err = bel.Render(ts, bel.GenerateAdditionalPreamble("// Note the alphabetical order which is unlike the \"natural\" order of the types and fields\n"))
	if err != nil {
		panic(err)
	}
}

func init() {
	examples["sort-alphabetically"] = SortAlphabetically
}
