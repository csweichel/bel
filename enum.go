package bel

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
)

// EnumHandler can determine if a type is an "enum" and retrieve its options
type EnumHandler interface {
	// IsEnum determines whether a type is an enum or not
	IsEnum(t reflect.Type) bool

	// GetMember returns all members of an enum
	GetMember(t reflect.Type) ([]TypescriptEnumMember, error)
}

// ParsedSourceEnumHandler discovers enums from type and const statements
type ParsedSourceEnumHandler struct {
	enums map[string][]TypescriptEnumMember
}

// NewParsedSourceEnumHandler creates a new enum handler that parses source code to discover enums
func NewParsedSourceEnumHandler(srcdir string) (*ParsedSourceEnumHandler, error) {
	fset := token.NewFileSet()
	pkgs := make(map[string]*ast.Package)
	err := filepath.Walk(srcdir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		ps, err := parser.ParseDir(fset, path, func(i os.FileInfo) bool { return true }, 0)
		if err != nil {
			return err
		}
		for n, pkg := range ps {
			pkgs[n] = pkg
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// the way the enum detection works at the moment this needs to be done in two passes
	enums := make(map[string][]TypescriptEnumMember)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, extractEnumTypes(enums))
		}
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, extractEnumValues(enums))
		}
	}

	return &ParsedSourceEnumHandler{enums: enums}, nil
}

func extractEnumTypes(enums map[string][]TypescriptEnumMember) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		if ts, ok := node.(*ast.TypeSpec); ok {
			enumName := ts.Name.Name
			if _, ok := ts.Type.(*ast.Ident); ok {
				enums[enumName] = make([]TypescriptEnumMember, 0)
			}

			return false
		}

		return true
	}
}

func extractEnumValues(enums map[string][]TypescriptEnumMember) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		if vs, ok := node.(*ast.ValueSpec); ok {
			if len(vs.Names) < 1 || len(vs.Values) < 1 {
				// TODO: add logging
				return false
			}

			var enumName string
			if tp, ok := vs.Type.(*ast.Ident); ok {
				enumName = tp.Name
			} else {
				return false
			}

			if members, ok := enums[enumName]; ok {
				name := vs.Names[0].Name
				value := vs.Values[0]
				if lit, ok := value.(*ast.BasicLit); ok {
					members = append(members, TypescriptEnumMember{
						Name:  name,
						Value: lit.Value,
					})
				}
				enums[enumName] = members
			}
			return false
		}

		return true
	}
}

// IsEnum returns true if the given type is an enumeration
func (h *ParsedSourceEnumHandler) IsEnum(t reflect.Type) bool {
	if _, ok := h.enums[t.Name()]; ok {
		return true
	}

	return false
}

// GetMember returns all members/values of an enum
func (h *ParsedSourceEnumHandler) GetMember(t reflect.Type) ([]TypescriptEnumMember, error) {
	if members, ok := h.enums[t.Name()]; ok {
		return members, nil
	}
	return nil, fmt.Errorf("no enum %s found", t.Name())
}
