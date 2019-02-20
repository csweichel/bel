package bel

type TypescriptKind string

const (
	TypescriptSimpleKind    TypescriptKind = "simple"
	TypescriptArrayKind     TypescriptKind = "array"
	TypescriptMapKind       TypescriptKind = "map"
	TypescriptInterfaceKind TypescriptKind = "iface"
)

type TypescriptType struct {
	Name    string
	Kind    TypescriptKind
	Members []TypescriptMember
	Params  []TypescriptType
}

type TypescriptMember struct {
	TypedElement
	IsOptional bool
	IsFunction bool
	Args       []TypedElement
}

type TypedElement struct {
	Name string
	Type TypescriptType
}
