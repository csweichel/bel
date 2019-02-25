package main

import (
	"github.com/32leaves/bel"
)

type StringEnum string

const (
	OptionOne   StringEnum = "option-one"
	OptionTwo   StringEnum = "option-two"
	OptionThree StringEnum = "option-three"
)

// For the enum detection to work the values need to be given explicitly.
// This UintEnum will not work with bel
type UintEnum uint32

const (
	OtherOptA UintEnum = iota
	OtherOptB
)

type SomeStruct struct {
	ThisOneWorks  StringEnum
	ThisOneDoesnt UintEnum
}

func ExtractEnums() {
	handler, err := bel.NewParsedSourceEnumHandler(".")
	if err != nil {
		panic(err)
	}

	ts, err := bel.Extract(SomeStruct{}, bel.WithEnumHandler(handler))
	if err != nil {
		panic(err)
	}

	err = bel.Render(ts)
	if err != nil {
		panic(err)
	}

	err = bel.Render(ts, bel.GenerateEnumAsSumType)
	if err != nil {
		panic(err)
	}
}

func init() {
	examples["enums"] = ExtractEnums
}
