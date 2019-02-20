package bel

import (
	"fmt"
	"reflect"
	"strings"
)

func Extract(s interface{}) (*TypescriptType, error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Struct {
		return extractStruct(t)
		// } else if t.Kind() == reflect.Interface {
		//     return extractInterface(t)
	} else {
		return nil, fmt.Errorf("cannot extract TS interface from %v", t.Kind())
	}
}

func extractStruct(t reflect.Type) (*TypescriptType, error) {
	fields := make([]TypescriptMember, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		exported := string(field.Name[0]) == strings.ToUpper(string(field.Name[0]))
		if !exported {
			continue
		}

		m, err := extractStructField(t.Field(i))
		if err != nil {
			return nil, err
		}
		fields = append(fields, *m)
	}

	return &TypescriptType{
		Name:    t.Name(),
		Kind:    TypescriptInterfaceKind,
		Members: fields,
	}, nil
}

func extractStructField(t reflect.StructField) (*TypescriptMember, error) {
	ttype := t.Type
	var tstype *TypescriptType
	if t.Type.Kind() == reflect.Ptr {
		ttype = ttype.Elem()
	}
	if ttype.Kind() == reflect.Struct {
		nme := ttype.Name()
		if nme == "" {
			res, err := extractStruct(ttype)
			if err != nil {
				return nil, err
			}
			tstype = res
		} else {
			tstype = &TypescriptType{Name: ttype.Name(), Kind: TypescriptSimpleKind}
		}
	} else {
		res, err := getPrimitiveType(ttype)
		if err != nil {
			return nil, err
		}
		tstype = res
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

func getPrimitiveType(t reflect.Type) (*TypescriptType, error) {
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
		elem, err := getPrimitiveType(t.Elem())
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
		key, err := getPrimitiveType(t.Key())
		if err != nil {
			return nil, err
		}
		elem, err := getPrimitiveType(t.Elem())
		if err != nil {
			return nil, err
		}
		return &TypescriptType{
			Kind:   TypescriptMapKind,
			Params: []TypescriptType{*key, *elem},
		}, nil
		// return mktype(fmt.Sprintf("{ [key: %s]: %s }", key, elem)), nil
	case reflect.Ptr:
		return getPrimitiveType(t.Elem())
	case reflect.String:
		return mktype("string"), nil
	}
	return nil, fmt.Errorf("cannot get primitive Typescript type for %v (%v)", t, kind)
}
