package main

// Generated with gopkg.in/reform.v1. Do not edit by hand.

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type sqliteMasterViewType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *sqliteMasterViewType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("sqlite_master").
func (v *sqliteMasterViewType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *sqliteMasterViewType) Columns() []string {
	return []string{"name"}
}

// NewStruct makes a new struct for that view or table.
func (v *sqliteMasterViewType) NewStruct() reform.Struct {
	return new(sqliteMaster)
}

// sqliteMasterView represents sqlite_master view or table in SQL database.
var sqliteMasterView = &sqliteMasterViewType{
	s: parse.StructInfo{Type: "sqliteMaster", SQLSchema: "", SQLName: "sqlite_master", Fields: []parse.FieldInfo{{Name: "Name", PKType: "", Column: "name"}}, PKFieldIndex: -1},
	z: new(sqliteMaster).Values(),
}

// String returns a string representation of this struct or record.
func (s sqliteMaster) String() string {
	res := make([]string, 1)
	res[0] = "Name: " + reform.Inspect(s.Name, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *sqliteMaster) Values() []interface{} {
	return []interface{}{
		s.Name,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *sqliteMaster) Pointers() []interface{} {
	return []interface{}{
		&s.Name,
	}
}

// View returns View object for that struct.
func (s *sqliteMaster) View() reform.View {
	return sqliteMasterView
}

// check interfaces
var (
	_ reform.View   = sqliteMasterView
	_ reform.Struct = (*sqliteMaster)(nil)
	_ fmt.Stringer  = (*sqliteMaster)(nil)
)

type sqliteTableInfoViewType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *sqliteTableInfoViewType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("dummy").
func (v *sqliteTableInfoViewType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *sqliteTableInfoViewType) Columns() []string {
	return []string{"cid", "name", "type", "notnull", "dflt_value", "pk"}
}

// NewStruct makes a new struct for that view or table.
func (v *sqliteTableInfoViewType) NewStruct() reform.Struct {
	return new(sqliteTableInfo)
}

// sqliteTableInfoView represents dummy view or table in SQL database.
var sqliteTableInfoView = &sqliteTableInfoViewType{
	s: parse.StructInfo{Type: "sqliteTableInfo", SQLSchema: "", SQLName: "dummy", Fields: []parse.FieldInfo{{Name: "CID", PKType: "", Column: "cid"}, {Name: "Name", PKType: "", Column: "name"}, {Name: "Type", PKType: "", Column: "type"}, {Name: "NotNull", PKType: "", Column: "notnull"}, {Name: "DefaultValue", PKType: "", Column: "dflt_value"}, {Name: "PK", PKType: "", Column: "pk"}}, PKFieldIndex: -1},
	z: new(sqliteTableInfo).Values(),
}

// String returns a string representation of this struct or record.
func (s sqliteTableInfo) String() string {
	res := make([]string, 6)
	res[0] = "CID: " + reform.Inspect(s.CID, true)
	res[1] = "Name: " + reform.Inspect(s.Name, true)
	res[2] = "Type: " + reform.Inspect(s.Type, true)
	res[3] = "NotNull: " + reform.Inspect(s.NotNull, true)
	res[4] = "DefaultValue: " + reform.Inspect(s.DefaultValue, true)
	res[5] = "PK: " + reform.Inspect(s.PK, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *sqliteTableInfo) Values() []interface{} {
	return []interface{}{
		s.CID,
		s.Name,
		s.Type,
		s.NotNull,
		s.DefaultValue,
		s.PK,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *sqliteTableInfo) Pointers() []interface{} {
	return []interface{}{
		&s.CID,
		&s.Name,
		&s.Type,
		&s.NotNull,
		&s.DefaultValue,
		&s.PK,
	}
}

// View returns View object for that struct.
func (s *sqliteTableInfo) View() reform.View {
	return sqliteTableInfoView
}

// check interfaces
var (
	_ reform.View   = sqliteTableInfoView
	_ reform.Struct = (*sqliteTableInfo)(nil)
	_ fmt.Stringer  = (*sqliteTableInfo)(nil)
)

func init() {
	parse.AssertUpToDate(&sqliteMasterView.s, new(sqliteMaster))
	parse.AssertUpToDate(&sqliteTableInfoView.s, new(sqliteTableInfo))
}
