package bel

import (
	"testing"

	"github.com/go-test/deep"
)

// InterfaceWithDocumentation has this documentation
type InterfaceWithDocumentation interface {
	// DoSomething also has documentation
	DoSomething(DoSomethingReq)
}

// DoSomethingReq is a struct with documentation
type DoSomethingReq struct {
	// Field has documentation as well
	Field string
}

func TestParsedSourceDocHandler(t *testing.T) {
	handler, err := NewParsedSourceDocHandler(".", "github.com/32leaves/")
	if err != nil {
		t.Error(err)
		return
	}

	extract, err := Extract((*InterfaceWithDocumentation)(nil), WithDocumentation(handler), FollowStructs)
	if err != nil {
		t.Error(err)
		return
	}

	expectation := []TypescriptType{
		{
			Name:    "DoSomethingReq",
			Comment: "DoSomethingReq is a struct with documentation",
			Kind:    TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "Field",
						Type: TypescriptType{
							Name: "string",
							Kind: TypescriptKind("simple"),
						},
					},
				},
			},
		},
		{
			Name:    "InterfaceWithDocumentation",
			Comment: "InterfaceWithDocumentation has this documentation",
			Kind:    TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "DoSomething",
						Type: TypescriptType{},
					},
					Comment:    "DoSomething also has documentation",
					IsFunction: true,
					Args: []TypedElement{
						{
							Name: "arg0",
							Type: TypescriptType{
								Name: "DoSomethingReq",
								Kind: TypescriptKind("simple"),
							},
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
