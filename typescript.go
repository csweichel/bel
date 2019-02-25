package bel

// TypescriptKind is the kind of a Typescript type (akin to the kind of Go types)
type TypescriptKind string

const (
	// TypescriptSimpleKind means the type is merely a symbol
	TypescriptSimpleKind TypescriptKind = "simple"
	// TypescriptArrayKind means the type is an array
	TypescriptArrayKind TypescriptKind = "array"
	// TypescriptMapKind means the type is a map/dict
	TypescriptMapKind TypescriptKind = "map"
	// TypescriptInterfaceKind means the type is an interface
	TypescriptInterfaceKind TypescriptKind = "iface"
	// TypescriptEnumKind means the type is an enum
	TypescriptEnumKind TypescriptKind = "enum"
)

// TypescriptType describes a type in the Typescript world
type TypescriptType struct {
	Name        string
	Kind        TypescriptKind
	Members     []TypescriptMember
	Params      []TypescriptType
	EnumMembers []TypescriptEnumMember
}

// TypescriptMember is a member of a Typescript interface
type TypescriptMember struct {
	TypedElement
	IsOptional bool
	IsFunction bool
	Args       []TypedElement
}

// TypescriptEnumMember is a member of a Typescript enum
type TypescriptEnumMember struct {
	Name  string
	Value string
}

// TypedElement pairs a name with a type
type TypedElement struct {
	Name string
	Type TypescriptType
}
