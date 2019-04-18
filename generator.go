package bel

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
)

const interfaceTemplate = `
{{ define "comment" -}}
{{- if .Comment }}
/**
 * {{ .Comment }}
 */
{{ end -}}
{{ end -}}
{{ define "iface" -}}
{
    {{ range .Members -}}
    {{- template "comment" . -}}
    {{ .Name }}{{ if .IsOptional }}?{{ end }}{{ if .IsFunction }}({{ template "args" . }}){{ end }}: {{ subt .Type | default "void" }}
    {{ end }}
}
{{ end -}}
{{- define "args" }}{{ range $idx, $val := .Args }}{{ if eq $idx 0 }}{{ else }}, {{ end }}{{ .Name }}: {{ subt .Type }}{{ end }}{{ end -}}
{{- define "simple" }}{{ .Name }}{{ end -}}
{{- define "map" }}{ [key: {{ subt (mapKeyType .) }}]: {{ subt (mapValType .) }} }{{ end -}}
{{- define "array" }}{{ subt (arrType .) }}[]{{ end -}}
{{- define "root-enum" }}{{- template "comment" . -}}export enum {{ .Name }} {
    {{ range .EnumMembers }}{{ .Name }} = {{ .Value }},
    {{ end }}
}{{ end -}}
{{- define "root-st-enum" }}{{- template "comment" . -}}export type {{ .Name }} =
    {{ range $idx, $val := .EnumMembers }}{{ if eq $idx 0 }}{{ else }} | {{ end }}{{ .Value }}{{ end }};
{{ end -}}
{{- define "root-iface" }}{{- template "comment" . -}}export interface {{ .Name }} {{ template "iface" . }} {{ end -}}
{{- .Preamble }}
{{ if .Namespace }}export namespace {{ .Namespace }} {
    {{ end -}}
{{- range .Types }}
{{ subtroot . }}
{{ end -}}
{{ if .Namespace }} } {{ end }}
`

type generateOptions struct {
	enumsAsSumTypes bool
	out             io.Writer
	Namespace       string
	Types           []TypescriptType
	Preamble        string
}

// GenerateOption is an option used with the Generate function
type GenerateOption func(*generateOptions)

// GenerateEnumAsSumType causes enums to be be rendered as sum types
func GenerateEnumAsSumType(opt *generateOptions) {
	opt.enumsAsSumTypes = true
}

// GenerateOutputTo sets the writer to which we'll write the generated TS code
func GenerateOutputTo(out io.Writer) GenerateOption {
	return func(opt *generateOptions) {
		opt.out = out
	}
}

// GenerateNamespace produces a namespace in which the generated types live
func GenerateNamespace(ns string) GenerateOption {
	return func(opt *generateOptions) {
		opt.Namespace = ns
	}
}

// GenerateAdditionalPreamble produces additional output at the beginning of the Typescript code
func GenerateAdditionalPreamble(preamble string) GenerateOption {
	return func(opt *generateOptions) {
		opt.Preamble += preamble
	}
}

// GeneratePreamble produces output at the beginning of the Typescript code
func GeneratePreamble(preamble string) GenerateOption {
	return func(opt *generateOptions) {
		opt.Preamble = preamble
	}
}

// Render produces TypeScript code
func Render(types []TypescriptType, cfg ...GenerateOption) error {
	opts := generateOptions{
		out:      os.Stdout,
		Preamble: fmt.Sprintf("// generated using github.com/32leaves/bel on %s\n// DO NOT MODIFY\n", time.Now()),
	}
	for _, c := range cfg {
		c(&opts)
	}

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
			if name == "" {
				return "", nil
			}

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
			}

			return "root-" + string(t.Kind)
		}),
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
	}

	var err error
	tpl, err = template.New("interface").Funcs(funcs).Parse(interfaceTemplate)
	if err != nil {
		return err
	}

	opts.Types = types

	r, w := io.Pipe()
	scanner := bufio.NewScanner(r)
	go func() {
		emptylines := 0
		for scanner.Scan() {
			line := scanner.Text()
			trimline := strings.TrimSpace(line)

			if trimline != "" {
				if emptylines > 1 {
					fmt.Fprintln(opts.out)
				}
				fmt.Fprintln(opts.out, line)
			}

			if trimline == "" {
				emptylines++
			} else {
				emptylines = 0
			}
		}
	}()

	return tpl.Execute(w, opts)
}
