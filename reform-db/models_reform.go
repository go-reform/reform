package main

// Generated with gopkg.in/reform.v1. Do not edit by hand.

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type tableViewType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("information_schema").
func (v *tableViewType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("tables").
func (v *tableViewType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *tableViewType) Columns() []string {
	return []string{"table_catalog", "table_schema", "table_name", "table_type"}
}

// NewStruct makes a new struct for that view or table.
func (v *tableViewType) NewStruct() reform.Struct {
	return new(table)
}

// tableView represents tables view or table in SQL database.
var tableView = &tableViewType{
	s: parse.StructInfo{Type: "table", SQLSchema: "information_schema", SQLName: "tables", Fields: []parse.FieldInfo{{Name: "Catalog", Type: "", Column: "table_catalog"}, {Name: "Schema", Type: "", Column: "table_schema"}, {Name: "Name", Type: "", Column: "table_name"}, {Name: "Type", Type: "", Column: "table_type"}}, PKFieldIndex: -1},
	z: new(table).Values(),
}

// String returns a string representation of this struct or record.
func (s table) String() string {
	res := make([]string, 4)
	res[0] = "Catalog: " + reform.Inspect(s.Catalog, true)
	res[1] = "Schema: " + reform.Inspect(s.Schema, true)
	res[2] = "Name: " + reform.Inspect(s.Name, true)
	res[3] = "Type: " + reform.Inspect(s.Type, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *table) Values() []interface{} {
	return []interface{}{
		s.Catalog,
		s.Schema,
		s.Name,
		s.Type,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *table) Pointers() []interface{} {
	return []interface{}{
		&s.Catalog,
		&s.Schema,
		&s.Name,
		&s.Type,
	}
}

// View returns View object for that struct.
func (s *table) View() reform.View {
	return tableView
}

// check interfaces
var (
	_ reform.View   = tableView
	_ reform.Struct = (*table)(nil)
	_ fmt.Stringer  = (*table)(nil)
)

type columnViewType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("information_schema").
func (v *columnViewType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("columns").
func (v *columnViewType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *columnViewType) Columns() []string {
	return []string{"table_catalog", "table_schema", "table_name", "column_name", "is_nullable", "data_type"}
}

// NewStruct makes a new struct for that view or table.
func (v *columnViewType) NewStruct() reform.Struct {
	return new(column)
}

// columnView represents columns view or table in SQL database.
var columnView = &columnViewType{
	s: parse.StructInfo{Type: "column", SQLSchema: "information_schema", SQLName: "columns", Fields: []parse.FieldInfo{{Name: "TableCatalog", Type: "", Column: "table_catalog"}, {Name: "TableSchema", Type: "", Column: "table_schema"}, {Name: "TableName", Type: "", Column: "table_name"}, {Name: "Name", Type: "", Column: "column_name"}, {Name: "IsNullable", Type: "", Column: "is_nullable"}, {Name: "Type", Type: "", Column: "data_type"}}, PKFieldIndex: -1},
	z: new(column).Values(),
}

// String returns a string representation of this struct or record.
func (s column) String() string {
	res := make([]string, 6)
	res[0] = "TableCatalog: " + reform.Inspect(s.TableCatalog, true)
	res[1] = "TableSchema: " + reform.Inspect(s.TableSchema, true)
	res[2] = "TableName: " + reform.Inspect(s.TableName, true)
	res[3] = "Name: " + reform.Inspect(s.Name, true)
	res[4] = "IsNullable: " + reform.Inspect(s.IsNullable, true)
	res[5] = "Type: " + reform.Inspect(s.Type, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *column) Values() []interface{} {
	return []interface{}{
		s.TableCatalog,
		s.TableSchema,
		s.TableName,
		s.Name,
		s.IsNullable,
		s.Type,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *column) Pointers() []interface{} {
	return []interface{}{
		&s.TableCatalog,
		&s.TableSchema,
		&s.TableName,
		&s.Name,
		&s.IsNullable,
		&s.Type,
	}
}

// View returns View object for that struct.
func (s *column) View() reform.View {
	return columnView
}

// check interfaces
var (
	_ reform.View   = columnView
	_ reform.Struct = (*column)(nil)
	_ fmt.Stringer  = (*column)(nil)
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
	s: parse.StructInfo{Type: "sqliteMaster", SQLSchema: "", SQLName: "sqlite_master", Fields: []parse.FieldInfo{{Name: "Name", Type: "", Column: "name"}}, PKFieldIndex: -1},
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
	s: parse.StructInfo{Type: "sqliteTableInfo", SQLSchema: "", SQLName: "dummy", Fields: []parse.FieldInfo{{Name: "CID", Type: "", Column: "cid"}, {Name: "Name", Type: "", Column: "name"}, {Name: "Type", Type: "", Column: "type"}, {Name: "NotNull", Type: "", Column: "notnull"}, {Name: "DefaultValue", Type: "", Column: "dflt_value"}, {Name: "PK", Type: "", Column: "pk"}}, PKFieldIndex: -1},
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
	parse.AssertUpToDate(&tableView.s, new(table))
	parse.AssertUpToDate(&columnView.s, new(column))
	parse.AssertUpToDate(&sqliteMasterView.s, new(sqliteMaster))
	parse.AssertUpToDate(&sqliteTableInfoView.s, new(sqliteTableInfo))
}
