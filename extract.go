package bel

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

type extractOption func(*Extractor)

// AnonStructNamer gives a name to an otherwise anonymous struct
type AnonStructNamer func(i reflect.StructField) string

// Namer translates Go type name convention to Typescript name convention.
// This function does not have to translate between Go types and Typescript types.
type Namer func(string) string

// Extractor pulls Typescript information from a Go structure
type Extractor struct {
	embedStructs    bool
	followStructs   bool
	noAnonStructs   bool
	anonStructNamer AnonStructNamer
	typeNamer       Namer
	enumHandler     EnumHandler

	result map[string]TypescriptType
}

var (
	// EmbedStructs produces a single monolithic structure where all
	// referenced/contained subtypes become a nested Typescript struct
	EmbedStructs extractOption = func(e *Extractor) {
		e.embedStructs = true
		e.followStructs = true
	}

	// FollowStructs enables transitive extraction of structs. By default
	// we just emit that struct's name.
	FollowStructs extractOption = func(e *Extractor) {
		e.followStructs = true
	}

	// NameAnonStructs enables non-monolithic extraction of anonymous structs.
	// Consider `struct { foo: struct { bar: int } }` where foo has an anonymous
	// struct as type - with NameAnonStructs set, we'd extract that struct as
	// its own Typescript interface.
	NameAnonStructs = func(namer AnonStructNamer) extractOption {
		return func(e *Extractor) {
			e.noAnonStructs = true
			e.anonStructNamer = namer
		}
	}

	// CustomNamer sets a custom function for translating Golang naming convention
	// to Typescript naming convention. This function does not have to translate
	// the type names, just the way they are written.
	CustomNamer = func(namer Namer) extractOption {
		return func(e *Extractor) {
			e.typeNamer = namer
		}
	}

	// WithEnumHandler configures an enum handler which detects and extracts enums from
	// types and constants.
	WithEnumHandler = func(handler EnumHandler) extractOption {
		return func(e *Extractor) {
			e.enumHandler = handler
		}
	}
)

// NewExtractor creates a new extractor
func NewExtractor(opts ...extractOption) *Extractor {
	result := &Extractor{
		embedStructs:  false,
		followStructs: false,
		typeNamer:     strcase.ToCamel,
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

func (e *Extractor) addResult(t *TypescriptType) {
	e.result[t.Name] = *t
}

func (e *Extractor) Extract(s interface{}) ([]TypescriptType, error) {
	e.result = make(map[string]TypescriptType)

	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Struct {
		estruct, err := e.extractStruct(t)
		if err != nil {
			return nil, err
		}
		e.addResult(estruct)
		// } else if t.Kind() == reflect.Interface {
		//     return extractInterface(t)
	} else {
		return nil, fmt.Errorf("cannot extract TS interface from %v", t.Kind())
	}

	res := make([]TypescriptType, 0)
	for _, e := range e.result {
		res = append(res, e)
	}
	return res, nil
}

func (e *Extractor) extractStruct(t reflect.Type) (*TypescriptType, error) {
	fields := make([]TypescriptMember, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		exported := string(field.Name[0]) == strings.ToUpper(string(field.Name[0]))
		if !exported {
			continue
		}

		m, err := e.extractStructField(t.Field(i))
		if err != nil {
			return nil, err
		}

		// skip fields named "-", see https://golang.org/pkg/encoding/json/#Marshal
		if m.Name != "-" {
			fields = append(fields, *m)
		}
	}

	return &TypescriptType{
		Name:    e.typeNamer(t.Name()),
		Kind:    TypescriptInterfaceKind,
		Members: fields,
	}, nil
}

func (e *Extractor) extractStructField(t reflect.StructField) (*TypescriptMember, error) {
	tstype, err := e.getType(t.Type, &t)
	if err != nil {
		return nil, err
	}

	optional := false
	name := t.Name
	if jsontag := t.Tag.Get("json"); jsontag != "" {
		segments := strings.Split(jsontag, ",")
		if len(segments) > 0 {
			segn := segments[0]
			if segn != "" {
				name = segn
			}
			segments = segments[1:]
		}
		for _, seg := range segments {
			if seg == "omitempty" {
				optional = true
			}
		}
	}

	return &TypescriptMember{
		TypedElement: TypedElement{
			Name: name,
			Type: *tstype,
		},
		IsOptional: optional,
		IsFunction: false,
	}, nil
}

func (e *Extractor) getType(ttype reflect.Type, t *reflect.StructField) (*TypescriptType, error) {
	var tstype *TypescriptType

	if ttype.Kind() == reflect.Ptr {
		ttype = ttype.Elem()
	}
	if ttype.Kind() == reflect.Struct {
		isanon := ttype.Name() == ""
		if isanon {
			astruct, err := e.extractStruct(ttype)
			if err != nil {
				return nil, err
			}

			if e.noAnonStructs {
				astructName := e.typeNamer(e.anonStructNamer(*t))
				astruct.Name = astructName
				e.addResult(astruct)
				tstype = &TypescriptType{Name: astructName, Kind: TypescriptSimpleKind}
			} else {
				tstype = astruct
			}
		} else if e.embedStructs {
			astruct, err := e.extractStruct(ttype)
			if err != nil {
				return nil, err
			}

			astruct.Name = ""
			tstype = astruct
		} else if e.followStructs {
			astruct, err := e.extractStruct(ttype)
			if err != nil {
				return nil, err
			}

			e.addResult(astruct)
			tstype = &TypescriptType{Name: astruct.Name, Kind: TypescriptSimpleKind}
		} else {
			tstype = &TypescriptType{Name: ttype.Name(), Kind: TypescriptSimpleKind}
		}
	} else if e.enumHandler != nil && e.enumHandler.IsEnum(ttype) {
		em, err := e.enumHandler.GetMember(ttype)
		if err != nil {
			return nil, err
		}
		enum := &TypescriptType{
			Name:        e.typeNamer(ttype.Name()),
			Kind:        TypescriptEnumKind,
			EnumMembers: em,
		}
		e.addResult(enum)
		tstype = &TypescriptType{Name: e.typeNamer(ttype.Name()), Kind: TypescriptSimpleKind}
	} else {
		res, err := e.getPrimitiveType(ttype)
		if err != nil {
			return nil, err
		}
		tstype = res
	}

	return tstype, nil
}

func (e *Extractor) getPrimitiveType(t reflect.Type) (*TypescriptType, error) {
	mktype := func(n string) *TypescriptType {
		return &TypescriptType{
			Kind: TypescriptSimpleKind,
			Name: n,
		}
	}

	kind := t.Kind()
	switch kind {
	case reflect.Bool:
		return mktype("boolean"), nil
	case reflect.Array,
		reflect.Slice:
		elem, err := e.getType(t.Elem(), nil)
		if err != nil {
			return nil, err
		}
		return &TypescriptType{
			Kind:   TypescriptArrayKind,
			Params: []TypescriptType{*elem},
		}, nil
	case reflect.Float32,
		reflect.Float64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return mktype("number"), nil
	case reflect.Map:
		key, err := e.getType(t.Key(), nil)
		if err != nil {
			return nil, err
		}
		elem, err := e.getType(t.Elem(), nil)
		if err != nil {
			return nil, err
		}
		return &TypescriptType{
			Kind:   TypescriptMapKind,
			Params: []TypescriptType{*key, *elem},
		}, nil
		// return mktype(fmt.Sprintf("{ [key: %s]: %s }", key, elem)), nil
	case reflect.Ptr:
		return e.getType(t.Elem(), nil)
	case reflect.String:
		return mktype("string"), nil
	}
	return nil, fmt.Errorf("cannot get primitive Typescript type for %v (%v)", t, kind)
}
