package bel

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
)

// ExtractOption is an option used with the Extract function
type ExtractOption func(*extractor)

// AnonStructNamer gives a name to an otherwise anonymous struct
type AnonStructNamer func(reflect.StructField) string

// TypeNamer translates Go type name convention to Typescript name convention.
// This function does not have to map Go types to Typescript types.
type TypeNamer func(reflect.Type) string

// extractor pulls Typescript information from a Go structure
type extractor struct {
	embedStructs    bool
	followStructs   bool
	noAnonStructs   bool
	sorter          func(a, b interface{}) bool
	anonStructNamer AnonStructNamer
	typeNamer       TypeNamer
	enumHandler     EnumHandler
	docHandler      DocHandler

	result map[string]TypescriptType
}

// EmbedStructs produces a single monolithic structure where all
// referenced/contained subtypes become a nested Typescript struct
func EmbedStructs(e *extractor) {
	e.embedStructs = true
	e.followStructs = true
}

// FollowStructs enables transitive extraction of structs. By default
// we just emit that struct's name.
func FollowStructs(e *extractor) {
	e.followStructs = true
}

// NameAnonStructs enables non-monolithic extraction of anonymous structs.
// Consider `struct { foo: struct { bar: int } }` where foo has an anonymous
// struct as type - with NameAnonStructs set, we'd extract that struct as
// its own Typescript interface.
func NameAnonStructs(namer AnonStructNamer) ExtractOption {
	return func(e *extractor) {
		e.noAnonStructs = true
		e.anonStructNamer = namer
	}
}

// CustomNamer sets a custom function for translating Golang naming convention
// to Typescript naming convention. This function does not have to translate
// the type names, just the way they are written.
func CustomNamer(namer TypeNamer) ExtractOption {
	return func(e *extractor) {
		e.typeNamer = namer
	}
}

// WithEnumerations configures an enum handler which detects and extracts enums from
// types and constants.
func WithEnumerations(handler EnumHandler) ExtractOption {
	return func(e *extractor) {
		e.enumHandler = handler
	}
}

// WithDocumentation configures a documentation handler which extracts documentation
// for types and methods.
func WithDocumentation(handler DocHandler) ExtractOption {
	return func(e *extractor) {
		e.docHandler = handler
	}
}

// SortAlphabetically sorts all types and their members alphabetically
func SortAlphabetically(e *extractor) {
	sorter := func(a, b interface{}) bool {
		ta, oka := a.(*TypescriptType)
		tb, okb := b.(*TypescriptType)
		if oka && okb {
			return ta.Name < tb.Name
		}

		tma, oka := a.(*TypescriptMember)
		tmb, okb := b.(*TypescriptMember)
		if oka && okb {
			return tma.Name < tmb.Name
		}

		tea, oka := a.(*TypescriptEnumMember)
		teb, okb := b.(*TypescriptEnumMember)
		if oka && okb {
			return tea.Name < teb.Name
		}

		return false
	}

	e.sorter = sorter
}

func (e *extractor) addResult(t *TypescriptType) {
	e.result[t.Name] = *t
}

// Extract uses reflection to extract the information required to generate Typescript code
func Extract(s interface{}, opts ...ExtractOption) ([]TypescriptType, error) {
	e := &extractor{
		embedStructs:  false,
		followStructs: false,
		typeNamer: func(t reflect.Type) string {
			return strcase.ToCamel(t.Name())
		},
		docHandler: (*nullDocHandler)(nil),
	}
	for _, opt := range opts {
		opt(e)
	}

	e.result = make(map[string]TypescriptType)

	t := reflect.TypeOf(s)
	if t == nil {
		return nil, fmt.Errorf("TypeOf(s) == nil")
	} else if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		estruct, err := e.extractStruct(t)
		if err != nil {
			return nil, err
		}
		e.addResult(estruct)
	} else if t.Kind() == reflect.Interface {
		et, err := e.extractInterface(t)
		if err != nil {
			return nil, err
		}
		e.addResult(et)
	} else {
		return nil, fmt.Errorf("cannot extract TS interface from %v", t.Kind())
	}

	res := make([]TypescriptType, 0)
	for _, e := range e.result {
		res = append(res, e)
	}
	if e.sorter != nil {
		sort.Slice(res, func(i, j int) bool {
			return e.sorter(&res[i], &res[j])
		})
	}
	return res, nil
}

func (e *extractor) extractInterface(t reflect.Type) (*TypescriptType, error) {
	if t.Kind() != reflect.Interface {
		return nil, fmt.Errorf("can only extract interface types")
	}

	methods := make([]TypescriptMember, t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		tm := t.Method(i)
		fnt := tm.Type

		var retval TypescriptType
		if fnt.NumOut() == 0 {
			// void
		} else if fnt.NumOut() <= 2 {
			errorInterface := reflect.TypeOf((*error)(nil)).Elem()
			if fnt.NumOut() == 2 && !fnt.Out(1).Implements(errorInterface) {
				return nil, fmt.Errorf("second return value must be an error in %s/%s", t.Name(), tm.Name)
			} else if fnt.Out(0).Implements(errorInterface) {
				// do not use this type - it's the error return
			} else {
				rv, err := e.getType(fnt.Out(0), nil)
				if err != nil {
					return nil, err
				}
				retval = *rv
			}
		} else {
			return nil, fmt.Errorf("cannot export more than two return values in %s/%s", t.Name(), tm.Name)
		}

		if fnt.IsVariadic() {
			return nil, fmt.Errorf("variadic functions are not supported: %s/%s", t.Name(), tm.Name)
		}
		args := make([]TypedElement, fnt.NumIn())
		for j := 0; j < fnt.NumIn(); j++ {
			at, err := e.getType(fnt.In(j), nil)
			if err != nil {
				return nil, err
			}
			args[j] = TypedElement{
				Name: fmt.Sprintf("arg%d", j),
				Type: *at,
			}
		}

		methods[i] = TypescriptMember{
			TypedElement: TypedElement{
				Name: tm.Name,
				Type: retval,
			},
			Comment:    e.docHandler.Method(t, tm),
			IsFunction: true,
			Args:       args,
		}
	}

	if e.sorter != nil {
		sort.Slice(methods, func(i, j int) bool {
			return e.sorter(&methods[i], &methods[j])
		})
	}
	res := &TypescriptType{
		Kind:    TypescriptInterfaceKind,
		Name:    e.typeNamer(t),
		Members: methods,
		Comment: e.docHandler.Type(t),
	}
	return res, nil
}

func (e *extractor) extractStruct(t reflect.Type) (*TypescriptType, error) {
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

	if e.sorter != nil {
		sort.Slice(fields, func(i, j int) bool {
			return e.sorter(&fields[i], &fields[j])
		})
	}
	return &TypescriptType{
		Name:    e.typeNamer(t),
		Comment: e.docHandler.Type(t),
		Kind:    TypescriptInterfaceKind,
		Members: fields,
	}, nil
}

func (e *extractor) extractStructField(t reflect.StructField) (*TypescriptMember, error) {
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

func (e *extractor) getType(ttype reflect.Type, t *reflect.StructField) (*TypescriptType, error) {
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
				astructName := e.anonStructNamer(*t)
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
			tstype = &TypescriptType{Name: e.typeNamer(ttype), Kind: TypescriptSimpleKind}
		}
	} else if e.enumHandler != nil && e.enumHandler.IsEnum(ttype) {
		em, err := e.enumHandler.GetMember(ttype)
		if err != nil {
			return nil, err
		}
		enum := &TypescriptType{
			Name:        e.typeNamer(ttype),
			Kind:        TypescriptEnumKind,
			EnumMembers: em,
		}
		e.addResult(enum)
		tstype = &TypescriptType{Name: e.typeNamer(ttype), Kind: TypescriptSimpleKind}
	} else {
		res, err := e.getPrimitiveType(ttype)
		if err != nil {
			return nil, err
		}
		tstype = res
	}

	return tstype, nil
}

func (e *extractor) getPrimitiveType(t reflect.Type) (*TypescriptType, error) {
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
	case reflect.Ptr:
		return e.getType(t.Elem(), nil)
	case reflect.String:
		return mktype("string"), nil
	}
	return nil, fmt.Errorf("cannot get primitive Typescript type for %v (%v)", t, kind)
}
