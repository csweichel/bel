package bel

import (
	"sort"
	"testing"

	_ "github.com/alecthomas/repr"

	"github.com/go-test/deep"
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
	sort.Slice(myenum, func(ia, ib int) bool { return myenum[ia].Name < myenum[ib].Name })

	expectation := []TypescriptEnumMember{
		{
			Name:  "MemberOne",
			Value: "\"member-one\"",
		},
		{
			Name:  "MemberThree",
			Value: "\"member-three\"",
		},
		{
			Name:  "MemberTwo",
			Value: "\"member-two\"",
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
	sort.Slice(enum, func(ia, ib int) bool { return enum[ia].Value < enum[ib].Value })

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

	extract, err := Extract(StructWithEnum{}, WithEnumHandler(handler))
	if err != nil {
		t.Error(err)
		return
	}
	sort.Slice(extract, func(ia, ib int) bool { return extract[ia].Name < extract[ib].Name })
	for i := range extract {
		sort.Slice(extract[i].Members, func(ia, ib int) bool { return extract[i].Members[ia].Name < extract[i].Members[ib].Name })
		sort.Slice(extract[i].EnumMembers, func(ia, ib int) bool { return extract[i].EnumMembers[ia].Value < extract[i].EnumMembers[ib].Value })
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
					Name:  "MemberThree",
					Value: "\"member-three\"",
				},
				{
					Name:  "MemberTwo",
					Value: "\"member-two\"",
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
				{
					TypedElement: TypedElement{
						Name: "Foo",
						Type: TypescriptType{
							Name: "MyEnum",
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
