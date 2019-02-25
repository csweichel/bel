package main

import (
	"fmt"
	"reflect"

	"github.com/32leaves/bel"
)

// NestedStuff contains a nested anonymous enum
type NestedStuff struct {
	SomeField struct {
		IAmNested string
		SoAmI     int32
	}
}

// NameAnonStructs demonstrates how to name anonymous structs
func NameAnonStructs() {
	anonNamer := func(t reflect.StructField) string {
		return fmt.Sprintf("WasAnon%s", t.Name)
	}
	ts, err := bel.Extract(NestedStuff{}, bel.NameAnonStructs(anonNamer))
	if err != nil {
		panic(err)
	}

	err = bel.Render(ts)
	if err != nil {
		panic(err)
	}
}

func init() {
	examples["name-anon-structs"] = NameAnonStructs
}
