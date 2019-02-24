package bel

import (
	"testing"
)

type DemoService interface {
	SayHello(name, msg string) (string, error)
}

func TestGenerateStuff(t *testing.T) {
	handler, err := NewParsedSourceEnumHandler(".")
	if err != nil {
		t.Error(err)
		return
	}

	extract, err := Extract((*DemoService)(nil),
		WithEnumHandler(handler),
		FollowStructs,
	)
	if err != nil {
		t.Error(err)
		return
	}

	err = Render(extract, GenerateEnumAsSumType, GenerateNamespace("foobar"))
	if err != nil {
		t.Error(err)
		return
	}
}
