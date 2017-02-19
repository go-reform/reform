package main

import (
	"text/template"

	"gopkg.in/reform.v1/parse"
)

type StructData struct {
	parse.StructInfo
	FieldComments []string
}

var (
	prologTemplate = template.Must(template.New("prolog").Parse(`
import (
	"time"
)
`))

	structTemplate = template.Must(template.New("struct").Parse(`
//go:generate reform

//reform:{{ .SQLName }}
type {{ .Type }} struct {
	{{- range $i, $f := .Fields }}
    {{ $f.Name }} {{ $f.Type }} ` + "`" + `reform:"{{ $f.Column }}{{ if eq $i $.PKFieldIndex }},pk{{ end }}"` + "`" + ` {{ index $.FieldComments $i }}
	{{- end }}
}
`))
)
