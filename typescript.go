package bel

type TypescriptKind string

const (
	TypescriptSimpleKind    TypescriptKind = "simple"
	TypescriptArrayKind     TypescriptKind = "array"
	TypescriptMapKind       TypescriptKind = "map"
	TypescriptInterfaceKind TypescriptKind = "iface"
	TypescriptEnumKind      TypescriptKind = "enum"
)

type TypescriptType struct {
	Name        string
	Kind        TypescriptKind
	Members     []TypescriptMember
	Params      []TypescriptType
	EnumMembers []TypescriptEnumMember
}

type TypescriptMember struct {
	TypedElement
	IsOptional bool
	IsFunction bool
	Args       []TypedElement
}

type TypescriptEnumMember struct {
	Name  string
	Value string
}

type TypedElement struct {
	Name string
	Type TypescriptType
}
