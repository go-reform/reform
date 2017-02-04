package main

import (
	"text/template"
)

var (
	prologTemplate = template.Must(template.New("prolog").Parse(`
import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)
`))

	// FIXME we abuse PKType name here, it should be renamed to Type

	structTemplate = template.Must(template.New("struct").Parse(`
//go:generate reform

//reform:{{ .SQLName }}
type {{ .Type }} struct {
	{{- range $i, $f := .Fields }}
    {{ $f.Name }} {{ $f.PKType }} ` + "`" + `reform:"{{ $f.Column }}"` + "`" + `
	{{- end }}
}`))
)
