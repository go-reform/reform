package parse_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/reform.v1/internal/test/models"
	"gopkg.in/reform.v1/internal/test/models/bogus"
	. "gopkg.in/reform.v1/parse"
)

//nolint:gochecknoglobals
var (
	person = StructInfo{
		Type:      "Person",
		SQLSchema: "",
		SQLName:   "people",
		Fields: []FieldInfo{
			{Name: "ID", Type: "int32", Column: "id"},
			{Name: "GroupID", Type: "*int32", Column: "group_id"},
			{Name: "Name", Type: "string", Column: "name"},
			{Name: "Email", Type: "*string", Column: "email"},
			{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
			{Name: "UpdatedAt", Type: "*time.Time", Column: "updated_at"},
		},
		PKFieldIndex: 0,
	}

	project = StructInfo{
		Type:      "Project",
		SQLSchema: "",
		SQLName:   "projects",
		Fields: []FieldInfo{
			{Name: "Name", Type: "string", Column: "name"},
			{Name: "ID", Type: "string", Column: "id"},
			{Name: "Start", Type: "time.Time", Column: "start"},
			{Name: "End", Type: "*time.Time", Column: "end"},
		},
		PKFieldIndex: 1,
	}

	personProject = StructInfo{
		Type:      "PersonProject",
		SQLSchema: "",
		SQLName:   "person_project",
		Fields: []FieldInfo{
			{Name: "PersonID", Type: "int32", Column: "person_id"},
			{Name: "ProjectID", Type: "string", Column: "project_id"},
		},
		PKFieldIndex: -1,
	}

	idOnly = StructInfo{
		Type:      "IDOnly",
		SQLSchema: "",
		SQLName:   "id_only",
		Fields: []FieldInfo{
			{Name: "ID", Type: "int32", Column: "id"},
		},
		PKFieldIndex: 0,
	}

	constraints = StructInfo{
		Type:      "Constraints",
		SQLSchema: "",
		SQLName:   "constraints",
		Fields: []FieldInfo{
			{Name: "I", Type: "int32", Column: "i"},
			{Name: "ID", Type: "string", Column: "id"},
		},
		PKFieldIndex: 1,
	}

	compositePk = StructInfo{
		Type:      "CompositePk",
		SQLSchema: "",
		SQLName:   "composite_pk",
		Fields: []FieldInfo{
			{Name: "I", Type: "int32", Column: "i"},
			{Name: "Name", Type: "string", Column: "name"},
			{Name: "J", Type: "string", Column: "j"},
		},
		PKFieldIndex: -1,
	}

	legacyPerson = StructInfo{
		Type:      "LegacyPerson",
		SQLSchema: "legacy",
		SQLName:   "people",
		Fields: []FieldInfo{
			{Name: "ID", Type: "int32", Column: "id"},
			{Name: "Name", Type: "*string", Column: "name"},
		},
		PKFieldIndex: 0,
	}

	extra = StructInfo{
		Type:      "Extra",
		SQLSchema: "",
		SQLName:   "extra",
		Fields: []FieldInfo{
			{Name: "ID", Type: "Integer", Column: "id"},
			{Name: "Name", Type: "*String", Column: "name"},
			{Name: "Byte", Type: "uint8", Column: "byte"},
			{Name: "Uint8", Type: "uint8", Column: "uint8"},
			{Name: "ByteP", Type: "*uint8", Column: "bytep"},
			{Name: "Uint8P", Type: "*uint8", Column: "uint8p"},
			{Name: "Bytes", Type: "[]uint8", Column: "bytes"},
			{Name: "Uint8s", Type: "[]uint8", Column: "uint8s"},
			{Name: "BytesA", Type: "[512]uint8", Column: "bytesa"},
			{Name: "Uint8sA", Type: "[512]uint8", Column: "uint8sa"},
			{Name: "BytesT", Type: "Bytes", Column: "bytest"},
			{Name: "Uint8sT", Type: "Uint8s", Column: "uint8st"},
		},
		PKFieldIndex: 0,
	}

	notExported = StructInfo{
		Type:      "notExported",
		SQLSchema: "",
		SQLName:   "not_exported",
		Fields: []FieldInfo{
			{Name: "ID", Type: "string", Column: "id"},
		},
		PKFieldIndex: 0,
	}
)

func TestFileGood(t *testing.T) {
	s, err := File(filepath.FromSlash("../internal/test/models/good.go"))
	assert.NoError(t, err)
	require.Len(t, s, 7)
	assert.Equal(t, person, s[0])
	assert.Equal(t, project, s[1])
	assert.Equal(t, personProject, s[2])
	assert.Equal(t, idOnly, s[3])
	assert.Equal(t, constraints, s[4])
	assert.Equal(t, compositePk, s[5])
	assert.Equal(t, legacyPerson, s[6])
}

func TestFileExtra(t *testing.T) {
	s, err := File(filepath.FromSlash("../internal/test/models/extra.go"))
	assert.NoError(t, err)
	require.Len(t, s, 2)
	assert.Equal(t, extra, s[0])
	assert.Equal(t, notExported, s[1])
}

func TestFileBogus(t *testing.T) {
	dir := filepath.FromSlash("../internal/test/models/bogus/")
	for file, msg := range map[string]error{
		"bogus1.go": errors.New(`reform: Bogus1 has anonymous field BogusType with "reform:" tag, it is not allowed`),
		"bogus2.go": errors.New(`reform: Bogus2 has anonymous field bogusType with "reform:" tag, it is not allowed`),
		"bogus3.go": errors.New(`reform: Bogus3 has non-exported field bogus with "reform:" tag, it is not allowed`),
		"bogus4.go": errors.New(`reform: Bogus4 has field Bogus with invalid "reform:" tag value, it is not allowed`),
		"bogus5.go": errors.New(`reform: Bogus5 has field Bogus with invalid "reform:" tag value, it is not allowed`),
		"bogus6.go": errors.New(`reform: Bogus6 has no fields with "reform:" tag, it is not allowed`),
		"bogus7.go": errors.New(`reform: Bogus7 has pointer field Bogus with with "pk" label in "reform:" tag, it is not allowed`),
		// "bogus8.go": errors.New(`reform: Bogus8 has pointer field Bogus with with "omitempty" label in "reform:" tag, it is not allowed`),
		"bogus8.go":  errors.New(`reform: Bogus8 has field Bogus with invalid "reform:" tag value, it is not allowed`),
		"bogus9.go":  errors.New(`reform: Bogus9 has field Bogus2 with "reform:" tag with duplicate column name bogus (used by Bogus1), it is not allowed`),
		"bogus10.go": errors.New(`reform: Bogus10 has field Bogus2 with with duplicate "pk" label in "reform:" tag (first used by Bogus1), it is not allowed`),
		"bogus11.go": errors.New(`reform: Bogus11 has slice field Bogus with with "pk" label in "reform:" tag, it is not allowed`),

		"bogus_ignore.go": nil,
	} {
		s, err := File(filepath.Join(dir, file))
		assert.Nil(t, s)
		assert.Equal(t, msg, err)
	}
}

func TestObjectGood(t *testing.T) {
	s, err := Object(new(models.Person), "", "people")
	assert.NoError(t, err)
	assert.Equal(t, &person, s)

	s, err = Object(new(models.Project), "", "projects")
	assert.NoError(t, err)
	assert.Equal(t, &project, s)

	s, err = Object(new(models.PersonProject), "", "person_project")
	assert.NoError(t, err)
	assert.Equal(t, &personProject, s)

	s, err = Object(new(models.IDOnly), "", "id_only")
	assert.NoError(t, err)
	assert.Equal(t, &idOnly, s)

	s, err = Object(new(models.Constraints), "", "constraints")
	assert.NoError(t, err)
	assert.Equal(t, &constraints, s)

	s, err = Object(new(models.CompositePk), "", "composite_pk")
	assert.NoError(t, err)
	assert.Equal(t, &compositePk, s)

	s, err = Object(new(models.LegacyPerson), "legacy", "people")
	assert.NoError(t, err)
	assert.Equal(t, &legacyPerson, s)
}

func TestObjectExtra(t *testing.T) {
	s, err := Object(new(models.Extra), "", "extra")
	assert.NoError(t, err)
	assert.Equal(t, &extra, s)

	// s, err := Object(new(models.notExported), "", "not_exported")
	// assert.NoError(t, err)
	// assert.Equal(t, &notExported, s)
}

func TestObjectBogus(t *testing.T) {
	for obj, msg := range map[interface{}]error{
		new(bogus.Bogus1): errors.New(`reform: Bogus1 has anonymous field BogusType with "reform:" tag, it is not allowed`),
		new(bogus.Bogus2): errors.New(`reform: Bogus2 has anonymous field bogusType with "reform:" tag, it is not allowed`),
		new(bogus.Bogus3): errors.New(`reform: Bogus3 has non-exported field bogus with "reform:" tag, it is not allowed`),
		new(bogus.Bogus4): errors.New(`reform: Bogus4 has field Bogus with invalid "reform:" tag value, it is not allowed`),
		new(bogus.Bogus5): errors.New(`reform: Bogus5 has field Bogus with invalid "reform:" tag value, it is not allowed`),
		new(bogus.Bogus6): errors.New(`reform: Bogus6 has no fields with "reform:" tag, it is not allowed`),
		new(bogus.Bogus7): errors.New(`reform: Bogus7 has pointer field Bogus with with "pk" label in "reform:" tag, it is not allowed`),
		// new(bogus.Bogus8): errors.New(`reform: Bogus8 has pointer field Bogus with with "omitempty" label in "reform:" tag, it is not allowed`),
		new(bogus.Bogus8):  errors.New(`reform: Bogus8 has field Bogus with invalid "reform:" tag value, it is not allowed`),
		new(bogus.Bogus9):  errors.New(`reform: Bogus9 has field Bogus2 with "reform:" tag with duplicate column name bogus (used by Bogus1), it is not allowed`),
		new(bogus.Bogus10): errors.New(`reform: Bogus10 has field Bogus2 with with duplicate "pk" label in "reform:" tag (first used by Bogus1), it is not allowed`),
		new(bogus.Bogus11): errors.New(`reform: Bogus11 has slice field Bogus with with "pk" label in "reform:" tag, it is not allowed`),

		// new(bogus.BogusIgnore): do not test,
	} {
		s, err := Object(obj, "", "bogus")
		assert.Nil(t, s)
		assert.Equal(t, msg, err)
	}
}

func TestHelpersGood(t *testing.T) {
	t.Parallel()

	t.Run("person", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, strings.TrimSpace(`
parse.StructInfo{
	Type: "Person",
	SQLName: "people",
	Fields: []parse.FieldInfo{
		{Name: "ID", Type: "int32", Column: "id"},
		{Name: "GroupID", Type: "*int32", Column: "group_id"},
		{Name: "Name", Type: "string", Column: "name"},
		{Name: "Email", Type: "*string", Column: "email"},
		{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
		{Name: "UpdatedAt", Type: "*time.Time", Column: "updated_at"},
	},
	PKFieldIndex: 0,
}`), person.GoString())
		assert.Equal(t, []string{"id", "group_id", "name", "email", "created_at", "updated_at"}, person.Columns())
		assert.Equal(t, strings.TrimSpace(`
[]string{
	"id",
	"group_id",
	"name",
	"email",
	"created_at",
	"updated_at",
}`), person.ColumnsGoString())
		assert.True(t, person.IsTable())
		assert.Equal(t, FieldInfo{Name: "ID", Type: "int32", Column: "id"}, person.PKField())
	})

	t.Run("project", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, strings.TrimSpace(`
parse.StructInfo{
	Type: "Project",
	SQLName: "projects",
	Fields: []parse.FieldInfo{
		{Name: "Name", Type: "string", Column: "name"},
		{Name: "ID", Type: "string", Column: "id"},
		{Name: "Start", Type: "time.Time", Column: "start"},
		{Name: "End", Type: "*time.Time", Column: "end"},
	},
	PKFieldIndex: 1,
}`), project.GoString())
		assert.Equal(t, []string{"name", "id", "start", "end"}, project.Columns())
		assert.Equal(t, strings.TrimSpace(`
[]string{
	"name",
	"id",
	"start",
	"end",
}`), project.ColumnsGoString())
		assert.True(t, project.IsTable())
		assert.Equal(t, FieldInfo{Name: "ID", Type: "string", Column: "id"}, project.PKField())
	})

	t.Run("personProject", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, strings.TrimSpace(`
parse.StructInfo{
	Type: "PersonProject",
	SQLName: "person_project",
	Fields: []parse.FieldInfo{
		{Name: "PersonID", Type: "int32", Column: "person_id"},
		{Name: "ProjectID", Type: "string", Column: "project_id"},
	},
	PKFieldIndex: -1,
}`), personProject.GoString())
		assert.Equal(t, []string{"person_id", "project_id"}, personProject.Columns())
		assert.Equal(t, strings.TrimSpace(`
[]string{
	"person_id",
	"project_id",
}`), personProject.ColumnsGoString())
		assert.False(t, personProject.IsTable())
	})

	t.Run("idOnly", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, strings.TrimSpace(`
parse.StructInfo{
	Type: "IDOnly",
	SQLName: "id_only",
	Fields: []parse.FieldInfo{
		{Name: "ID", Type: "int32", Column: "id"},
	},
	PKFieldIndex: 0,
}`), idOnly.GoString())
		assert.Equal(t, []string{"id"}, idOnly.Columns())
		assert.Equal(t, strings.TrimSpace(`
[]string{
	"id",
}`), idOnly.ColumnsGoString())
		assert.True(t, idOnly.IsTable())
		assert.Equal(t, FieldInfo{Name: "ID", Type: "int32", Column: "id"}, idOnly.PKField())
	})

	t.Run("constraints", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, strings.TrimSpace(`
parse.StructInfo{
	Type: "Constraints",
	SQLName: "constraints",
	Fields: []parse.FieldInfo{
		{Name: "I", Type: "int32", Column: "i"},
		{Name: "ID", Type: "string", Column: "id"},
	},
	PKFieldIndex: 1,
}`), constraints.GoString())
		assert.Equal(t, []string{"i", "id"}, constraints.Columns())
		assert.Equal(t, strings.TrimSpace(`
[]string{
	"i",
	"id",
}`), constraints.ColumnsGoString())
		assert.True(t, constraints.IsTable())
		assert.Equal(t, FieldInfo{Name: "ID", Type: "string", Column: "id"}, constraints.PKField())
	})

	t.Run("compositePk", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, strings.TrimSpace(`
parse.StructInfo{
	Type: "CompositePk",
	SQLName: "composite_pk",
	Fields: []parse.FieldInfo{
		{Name: "I", Type: "int32", Column: "i"},
		{Name: "Name", Type: "string", Column: "name"},
		{Name: "J", Type: "string", Column: "j"},
	},
	PKFieldIndex: -1,
}`), compositePk.GoString())
		assert.Equal(t, []string{"i", "name", "j"}, compositePk.Columns())
		assert.Equal(t, strings.TrimSpace(`
[]string{
	"i",
	"name",
	"j",
}`), compositePk.ColumnsGoString())
		assert.False(t, compositePk.IsTable())
	})

	t.Run("legacyPerson", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, strings.TrimSpace(`
parse.StructInfo{
	Type: "LegacyPerson",
	SQLSchema: "legacy",
	SQLName: "people",
	Fields: []parse.FieldInfo{
		{Name: "ID", Type: "int32", Column: "id"},
		{Name: "Name", Type: "*string", Column: "name"},
	},
	PKFieldIndex: 0,
}`), legacyPerson.GoString())
		assert.Equal(t, []string{"id", "name"}, legacyPerson.Columns())
		assert.Equal(t, strings.TrimSpace(`
[]string{
	"id",
	"name",
}`), legacyPerson.ColumnsGoString())
		assert.True(t, legacyPerson.IsTable())
		assert.Equal(t, FieldInfo{Name: "ID", Type: "int32", Column: "id"}, legacyPerson.PKField())
	})
}

func TestHelpersExtra(t *testing.T) {
	t.Run("extra", func(t *testing.T) {
		assert.Equal(t, strings.TrimSpace(`
parse.StructInfo{
	Type: "Extra",
	SQLName: "extra",
	Fields: []parse.FieldInfo{
		{Name: "ID", Type: "Integer", Column: "id"},
		{Name: "Name", Type: "*String", Column: "name"},
		{Name: "Byte", Type: "uint8", Column: "byte"},
		{Name: "Uint8", Type: "uint8", Column: "uint8"},
		{Name: "ByteP", Type: "*uint8", Column: "bytep"},
		{Name: "Uint8P", Type: "*uint8", Column: "uint8p"},
		{Name: "Bytes", Type: "[]uint8", Column: "bytes"},
		{Name: "Uint8s", Type: "[]uint8", Column: "uint8s"},
		{Name: "BytesA", Type: "[512]uint8", Column: "bytesa"},
		{Name: "Uint8sA", Type: "[512]uint8", Column: "uint8sa"},
		{Name: "BytesT", Type: "Bytes", Column: "bytest"},
		{Name: "Uint8sT", Type: "Uint8s", Column: "uint8st"},
	},
	PKFieldIndex: 0,
}`), extra.GoString())
		columns := []string{
			"id", "name",
			"byte", "uint8", "bytep", "uint8p", "bytes", "uint8s", "bytesa", "uint8sa", "bytest", "uint8st",
		}
		assert.Equal(t, columns, extra.Columns())
		columnsS := strings.TrimSpace(`
[]string{
	"id",
	"name",
	"byte",
	"uint8",
	"bytep",
	"uint8p",
	"bytes",
	"uint8s",
	"bytesa",
	"uint8sa",
	"bytest",
	"uint8st",
}`)
		assert.Equal(t, columnsS, extra.ColumnsGoString())
		assert.True(t, extra.IsTable())
		assert.Equal(t, FieldInfo{Name: "ID", Type: "Integer", Column: "id"}, extra.PKField())
	})

	t.Run("notExported", func(t *testing.T) {
		assert.Equal(t, strings.TrimSpace(`
parse.StructInfo{
	Type: "notExported",
	SQLName: "not_exported",
	Fields: []parse.FieldInfo{
		{Name: "ID", Type: "string", Column: "id"},
	},
	PKFieldIndex: 0,
}`), notExported.GoString())
		assert.Equal(t, []string{"id"}, notExported.Columns())
		assert.Equal(t, strings.TrimSpace(`
[]string{
	"id",
}`), notExported.ColumnsGoString())
		assert.True(t, notExported.IsTable())
		assert.Equal(t, FieldInfo{Name: "ID", Type: "string", Column: "id"}, notExported.PKField())
	})
}

func TestAssertUpToDate(t *testing.T) {
	AssertUpToDate(&person, new(models.Person))

	func() {
		defer func() {
			expected := `reform:
		Person struct information is not up-to-date.
		Typically this means that Person type definition was changed, but 'reform' command / 'go generate' was not run.

		`
			assert.Equal(t, expected, recover())
		}()

		p := person
		p.PKFieldIndex = 1
		AssertUpToDate(&p, new(models.Person))
	}()
}
