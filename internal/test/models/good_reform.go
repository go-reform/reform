package models

// generated with gopkg.in/reform.v1

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type personTable struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *personTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("people").
func (v *personTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *personTable) Columns() []string {
	return []string{"id", "group_id", "name", "email", "created_at", "updated_at"}
}

// NewStruct makes a new struct for that view or table.
func (v *personTable) NewStruct() reform.Struct {
	return new(Person)
}

// NewRecord makes a new record for that table.
func (v *personTable) NewRecord() reform.Record {
	return new(Person)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *personTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// PersonTable represents people view or table in SQL database.
var PersonTable = &personTable{
	s: parse.StructInfo{Type: "Person", SQLSchema: "", SQLName: "people", Fields: []parse.FieldInfo{{Name: "ID", Type: "int32", Column: "id"}, {Name: "GroupID", Type: "*int32", Column: "group_id"}, {Name: "Name", Type: "string", Column: "name"}, {Name: "Email", Type: "*string", Column: "email"}, {Name: "CreatedAt", Type: "time.Time", Column: "created_at"}, {Name: "UpdatedAt", Type: "*time.Time", Column: "updated_at"}}, PKFieldIndex: 0},
	z: new(Person).Values(),
}

// String returns a string representation of this struct or record.
func (s Person) String() string {
	res := make([]string, 6)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "GroupID: " + reform.Inspect(s.GroupID, true)
	res[2] = "Name: " + reform.Inspect(s.Name, true)
	res[3] = "Email: " + reform.Inspect(s.Email, true)
	res[4] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[5] = "UpdatedAt: " + reform.Inspect(s.UpdatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Person) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.GroupID,
		s.Name,
		s.Email,
		s.CreatedAt,
		s.UpdatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Person) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.GroupID,
		&s.Name,
		&s.Email,
		&s.CreatedAt,
		&s.UpdatedAt,
	}
}

// View returns View object for that struct.
func (s *Person) View() reform.View {
	return PersonTable
}

// Table returns Table object for that record.
func (s *Person) Table() reform.Table {
	return PersonTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Person) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Person) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Person) HasPK() bool {
	return s.ID != PersonTable.z[PersonTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Person) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = int32(i64)
	} else {
		s.ID = pk.(int32)
	}
}

// check interfaces
var (
	_ reform.View   = PersonTable
	_ reform.Struct = new(Person)
	_ reform.Table  = PersonTable
	_ reform.Record = new(Person)
	_ fmt.Stringer  = new(Person)
)

type projectTable struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *projectTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("projects").
func (v *projectTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *projectTable) Columns() []string {
	return []string{"name", "id", "start", "end"}
}

// NewStruct makes a new struct for that view or table.
func (v *projectTable) NewStruct() reform.Struct {
	return new(Project)
}

// NewRecord makes a new record for that table.
func (v *projectTable) NewRecord() reform.Record {
	return new(Project)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *projectTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// ProjectTable represents projects view or table in SQL database.
var ProjectTable = &projectTable{
	s: parse.StructInfo{Type: "Project", SQLSchema: "", SQLName: "projects", Fields: []parse.FieldInfo{{Name: "Name", Type: "string", Column: "name"}, {Name: "ID", Type: "string", Column: "id"}, {Name: "Start", Type: "time.Time", Column: "start"}, {Name: "End", Type: "*time.Time", Column: "end"}}, PKFieldIndex: 1},
	z: new(Project).Values(),
}

// String returns a string representation of this struct or record.
func (s Project) String() string {
	res := make([]string, 4)
	res[0] = "Name: " + reform.Inspect(s.Name, true)
	res[1] = "ID: " + reform.Inspect(s.ID, true)
	res[2] = "Start: " + reform.Inspect(s.Start, true)
	res[3] = "End: " + reform.Inspect(s.End, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Project) Values() []interface{} {
	return []interface{}{
		s.Name,
		s.ID,
		s.Start,
		s.End,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Project) Pointers() []interface{} {
	return []interface{}{
		&s.Name,
		&s.ID,
		&s.Start,
		&s.End,
	}
}

// View returns View object for that struct.
func (s *Project) View() reform.View {
	return ProjectTable
}

// Table returns Table object for that record.
func (s *Project) Table() reform.Table {
	return ProjectTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Project) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Project) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Project) HasPK() bool {
	return s.ID != ProjectTable.z[ProjectTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Project) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = string(i64)
	} else {
		s.ID = pk.(string)
	}
}

// check interfaces
var (
	_ reform.View   = ProjectTable
	_ reform.Struct = new(Project)
	_ reform.Table  = ProjectTable
	_ reform.Record = new(Project)
	_ fmt.Stringer  = new(Project)
)

type personProjectView struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *personProjectView) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("person_project").
func (v *personProjectView) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *personProjectView) Columns() []string {
	return []string{"person_id", "project_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *personProjectView) NewStruct() reform.Struct {
	return new(PersonProject)
}

// PersonProjectView represents person_project view or table in SQL database.
var PersonProjectView = &personProjectView{
	s: parse.StructInfo{Type: "PersonProject", SQLSchema: "", SQLName: "person_project", Fields: []parse.FieldInfo{{Name: "PersonID", Type: "int32", Column: "person_id"}, {Name: "ProjectID", Type: "string", Column: "project_id"}}, PKFieldIndex: -1},
	z: new(PersonProject).Values(),
}

// String returns a string representation of this struct or record.
func (s PersonProject) String() string {
	res := make([]string, 2)
	res[0] = "PersonID: " + reform.Inspect(s.PersonID, true)
	res[1] = "ProjectID: " + reform.Inspect(s.ProjectID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *PersonProject) Values() []interface{} {
	return []interface{}{
		s.PersonID,
		s.ProjectID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *PersonProject) Pointers() []interface{} {
	return []interface{}{
		&s.PersonID,
		&s.ProjectID,
	}
}

// View returns View object for that struct.
func (s *PersonProject) View() reform.View {
	return PersonProjectView
}

// check interfaces
var (
	_ reform.View   = PersonProjectView
	_ reform.Struct = new(PersonProject)
	_ fmt.Stringer  = new(PersonProject)
)

type legacyPersonTable struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("legacy").
func (v *legacyPersonTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("people").
func (v *legacyPersonTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *legacyPersonTable) Columns() []string {
	return []string{"id", "name"}
}

// NewStruct makes a new struct for that view or table.
func (v *legacyPersonTable) NewStruct() reform.Struct {
	return new(LegacyPerson)
}

// NewRecord makes a new record for that table.
func (v *legacyPersonTable) NewRecord() reform.Record {
	return new(LegacyPerson)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *legacyPersonTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// LegacyPersonTable represents people view or table in SQL database.
var LegacyPersonTable = &legacyPersonTable{
	s: parse.StructInfo{Type: "LegacyPerson", SQLSchema: "legacy", SQLName: "people", Fields: []parse.FieldInfo{{Name: "ID", Type: "int32", Column: "id"}, {Name: "Name", Type: "*string", Column: "name"}}, PKFieldIndex: 0},
	z: new(LegacyPerson).Values(),
}

// String returns a string representation of this struct or record.
func (s LegacyPerson) String() string {
	res := make([]string, 2)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Name: " + reform.Inspect(s.Name, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *LegacyPerson) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Name,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *LegacyPerson) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Name,
	}
}

// View returns View object for that struct.
func (s *LegacyPerson) View() reform.View {
	return LegacyPersonTable
}

// Table returns Table object for that record.
func (s *LegacyPerson) Table() reform.Table {
	return LegacyPersonTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *LegacyPerson) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *LegacyPerson) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *LegacyPerson) HasPK() bool {
	return s.ID != LegacyPersonTable.z[LegacyPersonTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *LegacyPerson) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = int32(i64)
	} else {
		s.ID = pk.(int32)
	}
}

// check interfaces
var (
	_ reform.View   = LegacyPersonTable
	_ reform.Struct = new(LegacyPerson)
	_ reform.Table  = LegacyPersonTable
	_ reform.Record = new(LegacyPerson)
	_ fmt.Stringer  = new(LegacyPerson)
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
	return []string{"id", "name", "bytes", "byte", "array"}
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
	s: parse.StructInfo{Type: "Extra", SQLSchema: "", SQLName: "extra", Fields: []parse.FieldInfo{{Name: "ID", Type: "Integer", Column: "id"}, {Name: "Name", Type: "*String", Column: "name"}, {Name: "Bytes", Type: "[]byte", Column: "bytes"}, {Name: "Byte", Type: "*byte", Column: "byte"}, {Name: "Array", Type: "[512]byte", Column: "array"}}, PKFieldIndex: 0},
	z: new(Extra).Values(),
}

// String returns a string representation of this struct or record.
func (s Extra) String() string {
	res := make([]string, 5)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Name: " + reform.Inspect(s.Name, true)
	res[2] = "Bytes: " + reform.Inspect(s.Bytes, true)
	res[3] = "Byte: " + reform.Inspect(s.Byte, true)
	res[4] = "Array: " + reform.Inspect(s.Array, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Extra) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Name,
		s.Bytes,
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
	parse.AssertUpToDate(&PersonTable.s, new(Person))
	parse.AssertUpToDate(&ProjectTable.s, new(Project))
	parse.AssertUpToDate(&PersonProjectView.s, new(PersonProject))
	parse.AssertUpToDate(&LegacyPersonTable.s, new(LegacyPerson))
	parse.AssertUpToDate(&ExtraTable.s, new(Extra))
}
