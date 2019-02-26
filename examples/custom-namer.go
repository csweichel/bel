package main

import (
	"github.com/32leaves/bel"
	"reflect"
)

// ThisStructsName is an empty struct
type ThisStructsName struct {
}

// CustomNamer demonstrates how to use custom type naming
func CustomNamer() {
	namer := func(t reflect.Type) string {
		return "Complete" + t.Name()
	}
	ts, err := bel.Extract(ThisStructsName{}, bel.CustomNamer(namer))
	if err != nil {
		panic(err)
	}

	err = bel.Render(ts)
	if err != nil {
		panic(err)
	}
}

func init() {
	examples["custom-namer"] = CustomNamer
}
