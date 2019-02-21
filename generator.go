package bel

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"
)

const interfaceTemplate = `
{{ define "iface" -}}
{
    {{ range .Members -}}
    {{ .Name }}{{ if .IsOptional }}?{{ end }}: {{ subt .Type }}
    {{ end }}
}
{{ end -}}
{{- define "simple" }}{{ .Name }}{{ end -}}
{{- define "map" }}{ [key: {{ subt (mapKeyType .) }}]: {{ subt (mapValType .) }} }{{ end -}}
{{- define "array" }}{{ subt (arrType .) }}[]{{ end -}}
{{- define "root-enum" }}export enum {{ .Name }} {
    {{ range .EnumMembers }}{{ .Name }} = {{ .Value }},
    {{ end }}
}{{ end -}}
{{- define "root-st-enum" }}export type {{ .Name }} =
    {{ range $idx, $val := .EnumMembers }}{{ if eq $idx 0 }}{{ else }} | {{ end }}{{ .Value }}{{ end }};
{{ end -}}
{{- define "root-iface" }} export interface {{ .Name }} {{ template "iface" . }} {{ end -}}
{{ subtroot . }}
`

type generateOptions struct {
	enumsAsSumTypes bool
	out             io.Writer
}

type generateOption func(*generateOptions)

var (
	// GenerateEnumAsSumType causes enums to be be rendered as sum types
	GenerateEnumAsSumType generateOption = func(opt *generateOptions) {
		opt.enumsAsSumTypes = true
	}

	// GenerateOutputTo sets the writer to which we'll write the generated TS code
	GenerateOutputTo = func(out io.Writer) generateOption {
		return func(opt *generateOptions) {
			opt.out = out
		}
	}
)

func Render(types []TypescriptType, cfg ...generateOption) error {
	opts := generateOptions{
		out: os.Stdout,
	}
	for _, c := range cfg {
		c(&opts)
	}

	for _, t := range types {
		err := opts.render(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (opts *generateOptions) render(ts TypescriptType) error {
	getParam := func(nme string, idx, minlen int) func(t TypescriptType) (*TypescriptType, error) {
		return func(t TypescriptType) (*TypescriptType, error) {
			if len(t.Params) != minlen {
				return nil, fmt.Errorf("map needs %d type params", minlen)
			}
			return &t.Params[idx], nil
		}
	}

	var tpl *template.Template
	executeTpl := func(selector func(t TypescriptType) string) func(t TypescriptType) (string, error) {
		if selector == nil {
			selector = func(t TypescriptType) string {
				return string(t.Kind)
			}
		}
		return func(t TypescriptType) (string, error) {
			name := selector(t)

			var b bytes.Buffer
			if err := tpl.ExecuteTemplate(&b, name, t); err != nil {
				return "", err
			}
			return b.String(), nil
		}
	}

	funcs := template.FuncMap{
		"mapKeyType": getParam("map", 0, 2),
		"mapValType": getParam("map", 1, 2),
		"arrType":    getParam("array", 0, 1),
		"subt":       executeTpl(nil),
		"subtroot": executeTpl(func(t TypescriptType) string {
			if t.Kind == TypescriptEnumKind && opts.enumsAsSumTypes {
				return "root-st-" + string(t.Kind)
			} else {
				return "root-" + string(t.Kind)
			}
		}),
	}

	var err error
	tpl, err = template.New("interface").Funcs(funcs).Parse(interfaceTemplate)
	if err != nil {
		return err
	}

	return tpl.Execute(opts.out, ts)
}
