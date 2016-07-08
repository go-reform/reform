package parse

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/reform.v1/internal/test/models"
	"gopkg.in/reform.v1/internal/test/models/bogus"
)

var (
	person = StructInfo{
		Type:    "Person",
		SQLName: "people",
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
		Type:    "Project",
		SQLName: "projects",
		Fields: []FieldInfo{
			{Name: "Name", Type: "string", Column: "name"},
			{Name: "ID", Type: "string", Column: "id"},
			{Name: "Start", Type: "time.Time", Column: "start"},
			{Name: "End", Type: "*time.Time", Column: "end"},
		},
		PKFieldIndex: 1,
	}

	personProject = StructInfo{
		Type:    "PersonProject",
		SQLName: "person_project",
		Fields: []FieldInfo{
			{Name: "PersonID", Type: "int32", Column: "person_id"},
			{Name: "ProjectID", Type: "string", Column: "project_id"},
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
		Type:    "Extra",
		SQLName: "extra",
		Fields: []FieldInfo{
			{Name: "ID", Type: "Integer", Column: "id"},
			{Name: "Name", Type: "*String", Column: "name"},
			{Name: "Bytes", Type: "[]byte", Column: "bytes"},
			{Name: "Byte", Type: "*byte", Column: "byte"},
			{Name: "Array", Type: "[512]byte", Column: "array"},
		},
		PKFieldIndex: 0,
	}
)

func TestFileGood(t *testing.T) {
	s, err := File("../internal/test/models/good.go")
	assert.NoError(t, err)
	require.Len(t, s, 5)
	assert.Equal(t, person, s[0])
	assert.Equal(t, project, s[1])
	assert.Equal(t, personProject, s[2])
	assert.Equal(t, legacyPerson, s[3])
	assert.Equal(t, extra, s[4])
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

	s, err = Object(new(models.LegacyPerson), "legacy", "people")
	assert.NoError(t, err)
	assert.Equal(t, &legacyPerson, s)

	s, err = Object(new(models.Extra), "", "extra")
	assert.NoError(t, err)
	assert.Equal(t, &extra, s)
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

		// new(bogus.BogusIgnore): do not test,
	} {
		s, err := Object(obj, "", "bogus")
		assert.Nil(t, s)
		assert.Equal(t, msg, err)
	}
}

func TestHelpers(t *testing.T) {
	assert.Equal(t, []string{"id", "group_id", "name", "email", "created_at", "updated_at"}, person.Columns())
	assert.True(t, person.IsTable())
	assert.Equal(t, FieldInfo{Name: "ID", Type: "int32", Column: "id"}, person.PKField())

	assert.Equal(t, []string{"name", "id", "start", "end"}, project.Columns())
	assert.True(t, project.IsTable())
	assert.Equal(t, FieldInfo{Name: "ID", Type: "string", Column: "id"}, project.PKField())

	assert.Equal(t, []string{"person_id", "project_id"}, personProject.Columns())
	assert.False(t, personProject.IsTable())
}
