package bel

import (
	"testing"
)

func TestGenerateStuff(t *testing.T) {
	handler, err := NewParsedSourceEnumHandler(".")
	if err != nil {
		t.Error(err)
		return
	}

	extract, err := NewExtractor(WithEnumHandler(handler)).Extract(StructWithEnum{})
	if err != nil {
		t.Error(err)
		return
	}

	err = Render(extract, GenerateEnumAsSumType)
	if err != nil {
		t.Error(err)
		return
	}
}
