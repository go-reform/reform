package models

// generated with gopkg.in/reform.v1

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type extraTable struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *extraTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("extra").
func (v *extraTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *extraTable) Columns() []string {
	return []string{"id", "name", "bytes", "bytes2", "byte", "array"}
}

// NewStruct makes a new struct for that view or table.
func (v *extraTable) NewStruct() reform.Struct {
	return new(Extra)
}

// NewRecord makes a new record for that table.
func (v *extraTable) NewRecord() reform.Record {
	return new(Extra)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *extraTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// ExtraTable represents extra view or table in SQL database.
var ExtraTable = &extraTable{
	s: parse.StructInfo{Type: "Extra", SQLSchema: "", SQLName: "extra", Fields: []parse.FieldInfo{{Name: "ID", PKType: "Integer", Column: "id"}, {Name: "Name", PKType: "", Column: "name"}, {Name: "Bytes", PKType: "", Column: "bytes"}, {Name: "Bytes2", PKType: "", Column: "bytes2"}, {Name: "Byte", PKType: "", Column: "byte"}, {Name: "Array", PKType: "", Column: "array"}}, PKFieldIndex: 0},
	z: new(Extra).Values(),
}

// String returns a string representation of this struct or record.
func (s Extra) String() string {
	res := make([]string, 6)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Name: " + reform.Inspect(s.Name, true)
	res[2] = "Bytes: " + reform.Inspect(s.Bytes, true)
	res[3] = "Bytes2: " + reform.Inspect(s.Bytes2, true)
	res[4] = "Byte: " + reform.Inspect(s.Byte, true)
	res[5] = "Array: " + reform.Inspect(s.Array, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Extra) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Name,
		s.Bytes,
		s.Bytes2,
		s.Byte,
		s.Array,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Extra) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Name,
		&s.Bytes,
		&s.Bytes2,
		&s.Byte,
		&s.Array,
	}
}

// View returns View object for that struct.
func (s *Extra) View() reform.View {
	return ExtraTable
}

// Table returns Table object for that record.
func (s *Extra) Table() reform.Table {
	return ExtraTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Extra) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Extra) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Extra) HasPK() bool {
	return s.ID != ExtraTable.z[ExtraTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Extra) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = Integer(i64)
	} else {
		s.ID = pk.(Integer)
	}
}

// check interfaces
var (
	_ reform.View   = ExtraTable
	_ reform.Struct = new(Extra)
	_ reform.Table  = ExtraTable
	_ reform.Record = new(Extra)
	_ fmt.Stringer  = new(Extra)
)

func init() {
	parse.AssertUpToDate(&ExtraTable.s, new(Extra))
}
