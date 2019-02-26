package bel

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
)

// DocHandler provides documentation for types
type DocHandler interface {
	// Type retrieves documentation for a type
	Type(t reflect.Type) string

	// Method retrieves documentation for an interface's method
	Method(parent reflect.Type, method reflect.Method) string
}

type nullDocHandler string

func (*nullDocHandler) Type(t reflect.Type) string {
	return ""
}

func (*nullDocHandler) Method(parent reflect.Type, method reflect.Method) string {
	return ""
}

// ParsedSourceDocHandler provides Go doc documentation from
type ParsedSourceDocHandler struct {
	pkgs map[string]*doc.Package
}

// NewParsedSourceDocHandler creates a new doc handler with a single pkg in its index
func NewParsedSourceDocHandler(srcdir, base string) (*ParsedSourceDocHandler, error) {
	res := &ParsedSourceDocHandler{pkgs: make(map[string]*doc.Package)}
	if err := res.AddToIndex(srcdir, base); err != nil {
		return nil, err
	}
	return res, nil
}

// AddToIndex adds another package to the handler's index. src is the path to the Go src folder of the package, pkg is its import path
func (h *ParsedSourceDocHandler) AddToIndex(src, pkg string) error {
	fset := token.NewFileSet()

	ps, err := parser.ParseDir(fset, src, func(i os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return err
	}
	for n, p := range ps {
		importPath := n
		if pkg != "" {
			importPath = fmt.Sprintf("%s/%s", strings.TrimRight(pkg, "/"), n)
		}
		h.pkgs[importPath] = doc.New(p, importPath, 0)
	}

	return nil
}

func (h *ParsedSourceDocHandler) findDoc(t reflect.Type) *doc.Type {
	pkg, ok := h.pkgs[t.PkgPath()]
	if !ok {
		return nil
	}

	for _, doct := range pkg.Types {
		if doct.Name == t.Name() {
			return doct
		}
	}

	return nil
}

// Type retrieves documentation for a type using the handler's index
func (h *ParsedSourceDocHandler) Type(t reflect.Type) string {
	doct := h.findDoc(t)
	if doct == nil {
		return ""
	}

	return strings.TrimSpace(doct.Doc)
}

// Method retrieves documentation for a method using the handler's index
func (h *ParsedSourceDocHandler) Method(parent reflect.Type, method reflect.Method) string {
	doct := h.findDoc(parent)
	if doct == nil {
		return ""
	}

	specs := doct.Decl.Specs
	if len(specs) < 1 {
		return ""
	}
	tspec, ok := specs[0].(*ast.TypeSpec)
	if !ok {
		return ""
	}

	ifspec, ok := tspec.Type.(*ast.InterfaceType)
	if !ok {
		return ""
	}

	for _, dm := range ifspec.Methods.List {
		if len(dm.Names) > 0 && dm.Names[0].Name == method.Name {
			return strings.TrimSpace(dm.Doc.Text())
		}
	}

	return ""
}
