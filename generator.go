package bel

import (
	"fmt"
	"io"
    "text/template"
    "bytes"
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
{{- define "map" }}{ [key: {{ subt (mapKeyType .) }}]: foo }{{ end -}}
{{- define "array" }}{{ subt (arrType .) }}[]{{ end -}}
export interface {{ .Name }} {{ template "iface" . }}
`

func (ts *TypescriptType) RenderInterface(w io.Writer) error {
    getParam := func(nme string, idx, minlen int) func(t TypescriptType) (*TypescriptType, error) {
        return func(t TypescriptType) (*TypescriptType, error) {
            if len(t.Params) != minlen {
                return nil, fmt.Errorf("map needs %d type params", minlen)
            }
            return &t.Params[idx], nil
        }
    }

    var tpl *template.Template
    funcs := template.FuncMap{
        "mapKeyType": getParam("map", 0, 2),
        "mapValType": getParam("map", 1, 2),
        "arrType": getParam("array", 0, 1),
        "subt": func(t TypescriptType) (string, error) {
            var b bytes.Buffer
			if err := tpl.ExecuteTemplate(&b, string(t.Kind), t); err != nil {
                return "", err
            }
            return b.String(), nil
        },
    }

    var err error
    tpl, err = template.New("interface").Funcs(funcs).Parse(interfaceTemplate)
    if err != nil {
        return err
    }

    return tpl.Execute(w, ts)
}
