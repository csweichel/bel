package bel

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

type MyTestStruct struct {
	notExported        string
	StringField        string
	OptionalField      string `json:",omitempty"`
	NamedField         int    `json:"thisFieldIsNamed"`
	NamedOptionalField int32  `json:"thisIsOptional,omitempty"`
	SkipThisField      string `json:"-"`
	Containment        AnotherTestStruct
	Referece           *AnotherTestStruct
}

type StructOfAllKind struct {
	BoolMember        bool
	StringArrayMember []string
	Float32Member     float32
	Float64Member     float64
	IntMember         int
	Int8Member        int8
	Int16Member       int16
	Int32Member       int32
	Int64Member       int64
	UintMember        uint
	Uint8Member       uint8
	Uint16Member      uint16
	Uint32Member      uint32
	Uint64Member      uint64
	MapMember         map[string]int
	PtrMember         *string
	StringMember      string
	AnonStruct        struct {
		AnonymousMember   string
		AnotherAnonMember int32
	}
}

type NestedStruct struct {
	Contains AnotherTestStruct
	Refers   *AnotherTestStruct
	Anon     struct {
		Foo string
	}
}

// AnotherTestStruct is just yet another struct
type AnotherTestStruct struct {
	// Foo has some documentation
	Foo string
	// Bar as well
	Bar bool
}

func TestExtractStruct(t *testing.T) {
	extract, err := Extract(MyTestStruct{})
	if err != nil {
		t.Error(err)
		return
	}

	// best generated with
	// repr.Print(extract)

	expectation := []TypescriptType{
		{
			Name: "MyTestStruct",
			Kind: TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "StringField",
						Type: TypescriptType{
							Name: "string",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "OptionalField",
						Type: TypescriptType{
							Name: "string",
							Kind: TypescriptKind("simple"),
						},
					},
					IsOptional: true,
				},
				{
					TypedElement: TypedElement{
						Name: "thisFieldIsNamed",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "thisIsOptional",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
					IsOptional: true,
				},
				{
					TypedElement: TypedElement{
						Name: "Containment",
						Type: TypescriptType{
							Name: "AnotherTestStruct",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Referece",
						Type: TypescriptType{
							Name: "AnotherTestStruct",
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

func TestNameAnonStructs(t *testing.T) {
	namer := func(t reflect.StructField) string {
		return t.Name
	}
	extract, err := Extract(NestedStruct{}, NameAnonStructs(namer))
	if err != nil {
		t.Error(err)
		return
	}

	// best generated with
	// repr.Print(extract)

	expectation := []TypescriptType{
		{
			Name: "Anon",
			Kind: TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "Foo",
						Type: TypescriptType{
							Name: "string",
							Kind: TypescriptKind("simple"),
						},
					},
				},
			},
		},
		{
			Name: "NestedStruct",
			Kind: TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "Contains",
						Type: TypescriptType{
							Name: "AnotherTestStruct",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Refers",
						Type: TypescriptType{
							Name: "AnotherTestStruct",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Anon",
						Type: TypescriptType{
							Name: "Anon",
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

func TestFollowStruct(t *testing.T) {
	extract, err := Extract(NestedStruct{}, FollowStructs)
	if err != nil {
		t.Error(err)
		return
	}

	// best generated with
	// repr.Print(extract)

	expectation := []TypescriptType{
		{
			Name: "AnotherTestStruct",
			Kind: TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "Foo",
						Type: TypescriptType{
							Name: "string",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Bar",
						Type: TypescriptType{
							Name: "boolean",
							Kind: TypescriptKind("simple"),
						},
					},
				},
			},
		},
		{
			Name: "NestedStruct",
			Kind: TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "Contains",
						Type: TypescriptType{
							Name: "AnotherTestStruct",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Refers",
						Type: TypescriptType{
							Name: "AnotherTestStruct",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Anon",
						Type: TypescriptType{
							Kind: TypescriptKind("iface"),
							Members: []TypescriptMember{
								{
									TypedElement: TypedElement{
										Name: "Foo",
										Type: TypescriptType{
											Name: "string",
											Kind: TypescriptKind("simple"),
										},
									},
								},
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

func TestEmbeddStruct(t *testing.T) {
	extract, err := Extract(NestedStruct{}, EmbedStructs)
	if err != nil {
		t.Error(err)
		return
	}

	// best generated with
	// repr.Print(extract)

	expectation := []TypescriptType{
		{
			Name: "NestedStruct",
			Kind: TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "Contains",
						Type: TypescriptType{
							Kind: TypescriptKind("iface"),
							Members: []TypescriptMember{
								{
									TypedElement: TypedElement{
										Name: "Foo",
										Type: TypescriptType{
											Name: "string",
											Kind: TypescriptKind("simple"),
										},
									},
								},
								{
									TypedElement: TypedElement{
										Name: "Bar",
										Type: TypescriptType{
											Name: "boolean",
											Kind: TypescriptKind("simple"),
										},
									},
								},
							},
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Refers",
						Type: TypescriptType{
							Kind: TypescriptKind("iface"),
							Members: []TypescriptMember{
								{
									TypedElement: TypedElement{
										Name: "Foo",
										Type: TypescriptType{
											Name: "string",
											Kind: TypescriptKind("simple"),
										},
									},
								},
								{
									TypedElement: TypedElement{
										Name: "Bar",
										Type: TypescriptType{
											Name: "boolean",
											Kind: TypescriptKind("simple"),
										},
									},
								},
							},
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Anon",
						Type: TypescriptType{
							Kind: TypescriptKind("iface"),
							Members: []TypescriptMember{
								{
									TypedElement: TypedElement{
										Name: "Foo",
										Type: TypescriptType{
											Name: "string",
											Kind: TypescriptKind("simple"),
										},
									},
								},
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

func TestStructOfAllKind(t *testing.T) {
	extract, err := Extract(StructOfAllKind{})
	if err != nil {
		t.Error(err)
		return
	}

	expectation := []TypescriptType{
		{
			Name: "StructOfAllKind",
			Kind: TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "BoolMember",
						Type: TypescriptType{
							Name: "boolean",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "StringArrayMember",
						Type: TypescriptType{
							Kind: TypescriptKind("array"),
							Params: []TypescriptType{
								{
									Name: "string",
									Kind: TypescriptKind("simple"),
								},
							},
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Float32Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Float64Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "IntMember",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Int8Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Int16Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Int32Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Int64Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "UintMember",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Uint8Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Uint16Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Uint32Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "Uint64Member",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "MapMember",
						Type: TypescriptType{
							Kind: TypescriptKind("map"),
							Params: []TypescriptType{
								{
									Name: "string",
									Kind: TypescriptKind("simple"),
								},
								{
									Name: "number",
									Kind: TypescriptKind("simple"),
								},
							},
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "PtrMember",
						Type: TypescriptType{
							Name: "string",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "StringMember",
						Type: TypescriptType{
							Name: "string",
							Kind: TypescriptKind("simple"),
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "AnonStruct",
						Type: TypescriptType{
							Kind: TypescriptKind("iface"),
							Members: []TypescriptMember{
								{
									TypedElement: TypedElement{
										Name: "AnonymousMember",
										Type: TypescriptType{
											Name: "string",
											Kind: TypescriptKind("simple"),
										},
									},
								},
								{
									TypedElement: TypedElement{
										Name: "AnotherAnonMember",
										Type: TypescriptType{
											Name: "number",
											Kind: TypescriptKind("simple"),
										},
									},
								},
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

type MyInterface interface {
	FirstOp(arg MyTestStruct) (int, error)
	SecondOp(arg0 int32, arg1 *StructOfAllKind) (*StructOfAllKind, error)
	VoidOp()
}

func TestExtractInterface(t *testing.T) {
	extract, _ := Extract((*MyInterface)(nil))

	// best generated using
	// repr.Print(extract)
	expectation := []TypescriptType{
		{
			Name: "MyInterface",
			Kind: TypescriptKind("iface"),
			Members: []TypescriptMember{
				{
					TypedElement: TypedElement{
						Name: "FirstOp",
						Type: TypescriptType{
							Name: "number",
							Kind: TypescriptKind("simple"),
						},
					},
					IsFunction: true,
					Args: []TypedElement{
						{
							Name: "arg0",
							Type: TypescriptType{
								Name: "MyTestStruct",
								Kind: TypescriptKind("simple"),
							},
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "SecondOp",
						Type: TypescriptType{
							Name: "StructOfAllKind",
							Kind: TypescriptKind("simple"),
						},
					},
					IsFunction: true,
					Args: []TypedElement{
						{
							Name: "arg0",
							Type: TypescriptType{
								Name: "number",
								Kind: TypescriptKind("simple"),
							},
						},
						{
							Name: "arg1",
							Type: TypescriptType{
								Name: "StructOfAllKind",
								Kind: TypescriptKind("simple"),
							},
						},
					},
				},
				{
					TypedElement: TypedElement{
						Name: "VoidOp",
						Type: TypescriptType{},
					},
					IsFunction: true,
					Args:       []TypedElement{},
				},
			},
		},
	}
	diff := deep.Equal(expectation, extract)
	for _, d := range diff {
		t.Error(d)
	}
}
