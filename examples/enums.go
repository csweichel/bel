package main

import (
	"github.com/32leaves/bel"
)

// StringEnum is an enumeration with string values
type StringEnum string

const (
	// OptionOne is an enum value
	OptionOne StringEnum = "option-one"
	// OptionTwo is an enum value
	OptionTwo StringEnum = "option-two"
	// OptionThree is an enum value
	OptionThree StringEnum = "option-three"
)

// UintEnum demonstrates which kind of enums do not work
// For the enum detection to work the values need to be given explicitly.
// This UintEnum will not work with bel
type UintEnum uint32

const (
	// OtherOptA is an enum value
	OtherOptA UintEnum = iota
	// OtherOptB is an enum value
	OtherOptB
)

// SomeStruct uses enumerations
type SomeStruct struct {
	ThisOneWorks  StringEnum
	ThisOneDoesnt UintEnum
}

// ExtractEnums demonstrates how to use an enum handler
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
