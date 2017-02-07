package main

import (
	"text/template"
)

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
    {{ $f.Name }} {{ $f.Type }} ` + "`" + `reform:"{{ $f.Column }}"` + "`" + `
	{{- end }}
}`))
)
