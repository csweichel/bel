package bel

import (
	// "github.com/alecthomas/repr"
	"github.com/go-test/deep"
	"testing"
)

type MyEnum string

const (
	MemberOne   MyEnum = "member-one"
	MemberTwo   MyEnum = "member-two"
	MemberThree MyEnum = "member-three"
)

type MyOtherEnum int

const (
	OtherEnumOne   MyOtherEnum = 0
	OtherEnumTwo   MyOtherEnum = 1
	OtherEnumThree MyOtherEnum = 2
	OtherEnumFour  MyOtherEnum = 3
)

type StructWithEnum struct {
    Foo MyEnum
    Bar MyOtherEnum
    Baz string
}

func TestParseStringEnum(t *testing.T) {
	handler, err := NewParsedSourceEnumHandler(".")
	if err != nil {
		t.Error(err)
		return
	}

	myenum, ok := handler.enums["MyEnum"]
	if !ok {
		t.Errorf("did not find MyEnum enum in sources")
		return
	}

	expectation := []TypescriptEnumMember{
		{
			Name:  "MemberOne",
			Value: "\"member-one\"",
		},
		{
			Name:  "MemberTwo",
			Value: "\"member-two\"",
		},
		{
			Name:  "MemberThree",
			Value: "\"member-three\"",
		},
	}
	diff := deep.Equal(expectation, myenum)
	for _, d := range diff {
		t.Error(d)
	}
}

func TestParseIntEnum(t *testing.T) {
	handler, err := NewParsedSourceEnumHandler(".")
	if err != nil {
		t.Error(err)
		return
	}

	enum, ok := handler.enums["MyOtherEnum"]
	if !ok {
		t.Errorf("did not find MyOtherEnum enum in sources")
		return
	}

	expectation := []TypescriptEnumMember{
		{
			Name:  "OtherEnumOne",
			Value: "0",
		},
		{
			Name:  "OtherEnumTwo",
			Value: "1",
		},
		{
			Name:  "OtherEnumThree",
			Value: "2",
		},
		{
			Name:  "OtherEnumFour",
			Value: "3",
		},
	}
	diff := deep.Equal(expectation, enum)
	for _, d := range diff {
		t.Error(d)
	}
}

func TestExtractIntEnum(t *testing.T) {
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

	expectation := []TypescriptType{
		{
			Name: "MyEnum",
			Kind: TypescriptKind("enum"),
			EnumMembers: []TypescriptEnumMember{
				{
					Name:  "MemberOne",
					Value: "\"member-one\"",
				},
				{
					Name:  "MemberTwo",
					Value: "\"member-two\"",
				},
				{
					Name:  "MemberThree",
					Value: "\"member-three\"",
				},
			},
		},
		{
			Name: "MyOtherEnum",
			Kind: TypescriptKind("enum"),
			EnumMembers: []TypescriptEnumMember{
				{
					Name:  "OtherEnumOne",
					Value: "0",
				},
				{
					Name:  "OtherEnumTwo",
					Value: "1",
				},
				{
					Name:  "OtherEnumThree",
					Value: "2",
				},
				{
					Name:  "OtherEnumFour",
					Value: "3",
				},
			},
		},
		{
			Name: "StructWithEnum",
			Kind: TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "Foo",
						Type: TypescriptType{
							Name: "MyEnum",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Bar",
						Type: TypescriptType{
							Name: "MyOtherEnum",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Baz",
						Type: TypescriptType{
							Name: "string",
							Kind: TypescriptKind("simple"),
						},
					},
				},
			},
		},
	}
	diff := deep.Equal(expectation, extract)
	for _, d := range diff {
		t.Error(d)
	}
}
