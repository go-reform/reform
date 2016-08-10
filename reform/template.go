package main

import (
	"text/template"

	"gopkg.in/reform.v1/parse"
)

// StructData represents struct info for XXX_reform.go file generation.
type StructData struct {
	parse.StructInfo
	TableType string
	TableVar  string
}

var (
	prologTemplate = template.Must(template.New("prolog").Parse(`
// generated with gopkg.in/reform.v1

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)
`))

	structTemplate = template.Must(template.New("struct").Parse(`
type {{ .TableType }} struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("{{ .SQLSchema }}").
func (v *{{ .TableType }}) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("{{ .SQLName }}").
func (v *{{ .TableType }}) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *{{ .TableType }}) Columns() []string {
	return {{ printf "%#v" .Columns }}
}

// NewStruct makes a new struct for that view or table.
func (v *{{ .TableType }}) NewStruct() reform.Struct {
	return new({{ .Type }})
}

{{- if .IsTable }}

// NewRecord makes a new record for that table.
func (v *{{ .TableType }}) NewRecord() reform.Record {
	return new({{ .Type }})
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *{{ .TableType }}) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

{{- end }}

// {{ .TableVar }} represents {{ .SQLName }} view or table in SQL database.
var {{ .TableVar }} = &{{ .TableType }} {
	s: {{ printf "%#v" .StructInfo }},
	z: new({{ .Type }}).Values(),
}

// String returns a string representation of this struct or record.
func (s {{ .Type }}) String() string {
	res := make([]string, {{ len .Fields }})
	{{- range $i, $f := .Fields }}
	res[{{ $i }}] = "{{ $f.Name }}: " + reform.Inspect(s.{{ $f.Name }}, true)
	{{- end }}
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *{{ .Type }}) Values() []interface{} {
	return []interface{}{ {{- range .Fields }}
		s.{{ .Name }}, {{- end }}
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *{{ .Type }}) Pointers() []interface{} {
	return []interface{}{ {{- range .Fields }}
		&s.{{ .Name }}, {{- end }}
	}
}

// View returns View object for that struct.
func (s *{{ .Type }}) View() reform.View {
	return {{ .TableVar }}
}

{{- if .IsTable }}

// Table returns Table object for that record.
func (s *{{ .Type }}) Table() reform.Table {
	return {{ .TableVar }}
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *{{ .Type }}) PKValue() interface{} {
	return s.{{ .PKField.Name }}
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *{{ .Type }}) PKPointer() interface{} {
	return &s.{{ .PKField.Name }}
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *{{ .Type }}) HasPK() bool {
	return s.{{ .PKField.Name }} != {{ .TableVar }}.z[{{ .TableVar }}.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *{{ .Type }}) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.{{ .PKField.Name }} = {{ .PKField.PKType }}(i64)
	} else {
		s.{{ .PKField.Name }} = pk.({{ .PKField.PKType }})
	}
}

{{- end }}

// check interfaces
var (
	_ reform.View   = {{ .TableVar }}
	_ reform.Struct = new({{ .Type }})
{{- if .IsTable }}
	_ reform.Table  = {{ .TableVar }}
	_ reform.Record = new({{ .Type }})
{{- end }}
	_ fmt.Stringer  = new({{ .Type }})
)
`))

	initTemplate = template.Must(template.New("init").Parse(`
func init() {
	{{- range $i, $sd := . }}
	parse.AssertUpToDate(&{{ $sd.TableVar }}.s, new({{ $sd.Type }}))
	{{- end }}
}
`))
)
