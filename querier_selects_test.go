package reform_test

import (
	"time"

	"github.com/AlekSi/pointer"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	. "gopkg.in/reform.v1/internal/test/models"
)

var (
	goCreated     = time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
	personCreated = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	baronStart    = time.Date(2014, 6, 1, 0, 0, 0, 0, time.UTC)
	baronEnd      = time.Date(2016, 2, 21, 0, 0, 0, 0, time.UTC)
	queenStart    = time.Date(2016, 1, 15, 0, 0, 0, 0, time.UTC)
)

func (s *ReformSuite) TestSelectOneTo() {
	var person Person
	err := s.q.SelectOneTo(&person, "WHERE id = "+s.q.Placeholder(1), 1)
	s.NoError(err)
	s.Equal(Person{ID: 1, GroupID: pointer.ToInt32(65534), Name: "Denis Mills", CreatedAt: goCreated}, person)

	var project Project
	err = s.q.SelectOneTo(&project, "WHERE id = "+s.q.Placeholder(1), "baron")
	s.NoError(err)
	expected := Project{ID: "baron", Name: "Vicious Baron", Start: baronStart, End: &baronEnd}
	s.Equal(expected, project)

	err = s.q.SelectOneTo(&project, "WHERE id IS NULL")
	s.Equal(expected, project) // expect old value
	s.Equal(reform.ErrNoRows, err)

	err = s.q.SelectOneTo(&project, "WHERE invalid_tail")
	s.Equal(expected, project) // expect old value
	s.Error(err)
	s.NotEqual(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestSelectOneFrom() {
	person, err := s.q.SelectOneFrom(PersonTable, "WHERE id = "+s.q.Placeholder(1), 1)
	s.NoError(err)
	s.Equal(&Person{ID: 1, GroupID: pointer.ToInt32(65534), Name: "Denis Mills", CreatedAt: goCreated}, person)

	project, err := s.q.SelectOneFrom(ProjectTable, "WHERE id = "+s.q.Placeholder(1), "baron")
	s.NoError(err)
	s.Equal(&Project{ID: "baron", Name: "Vicious Baron", Start: baronStart, End: &baronEnd}, project)

	project, err = s.q.SelectOneFrom(ProjectTable, "WHERE id IS NULL")
	s.Nil(project)
	s.Equal(reform.ErrNoRows, err)

	project, err = s.q.SelectOneFrom(ProjectTable, "WHERE invalid_tail")
	s.Nil(project)
	s.Error(err)
}

func (s *ReformSuite) TestSelectRows() {
	rows, err := s.q.SelectRows(PersonTable, "WHERE name = "+s.q.Placeholder(1)+" ORDER BY id", "Elfrieda Abbott")
	s.NotNil(rows)
	s.NoError(err)
	defer rows.Close()

	var person Person
	err = s.q.NextRow(&person, rows)
	s.NoError(err)
	expected := Person{ID: 102, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", Email: pointer.ToString("elfrieda_abbott@example.org"), CreatedAt: personCreated}
	s.Equal(expected, person)

	err = s.q.NextRow(&person, rows)
	s.NoError(err)
	expected = Person{ID: 103, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", CreatedAt: personCreated}
	s.Equal(expected, person)

	err = s.q.NextRow(&person, rows)
	s.Equal(err, reform.ErrNoRows)
	s.Equal(expected, person) // expect old value

	rows, err = s.q.SelectRows(ProjectTable, "WHERE id IS NULL")
	s.NotNil(rows)
	s.NoError(err)
	defer rows.Close()

	var project Project
	err = s.q.NextRow(&project, rows)
	s.Equal(reform.ErrNoRows, err)
	s.Equal(Project{}, project)

	rows, err = s.q.SelectRows(ProjectTable, "WHERE invalid_tail")
	s.Nil(rows)
	s.Error(err)
	s.NotEqual(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestSelectAllFrom() {
	structs, err := s.q.SelectAllFrom(PersonTable, "WHERE name = "+s.q.Placeholder(1)+" ORDER BY id", "Elfrieda Abbott")
	s.NoError(err)
	s.Len(structs, 2)
	s.Equal([]reform.Struct{
		&Person{ID: 102, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", Email: pointer.ToString("elfrieda_abbott@example.org"), CreatedAt: personCreated},
		&Person{ID: 103, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", CreatedAt: personCreated},
	}, structs)

	structs, err = s.q.SelectAllFrom(ProjectTable, "WHERE id IS NULL")
	s.Nil(structs)
	s.NoError(err)

	structs, err = s.q.SelectAllFrom(ProjectTable, "WHERE invalid_tail")
	s.Nil(structs)
	s.Error(err)
	s.NotEqual(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestFindOneTo() {
	var person Person
	err := s.q.FindOneTo(&person, "id", 102)
	s.NoError(err)
	s.Equal(
		Person{ID: 102, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", Email: pointer.ToString("elfrieda_abbott@example.org"), CreatedAt: personCreated},
		person,
	)

	var project Project
	err = s.q.FindOneTo(&project, "id", "queen")
	s.NoError(err)
	expected := Project{ID: "queen", Name: "Thirsty Queen", Start: queenStart}
	s.Equal(expected, project)

	err = s.q.FindOneTo(&project, "id", nil)
	s.Equal(expected, project) // expect old value
	s.Equal(reform.ErrNoRows, err)

	err = s.q.FindOneTo(&project, "invalid_column", nil)
	s.Equal(expected, project) // expect old value
	s.Error(err)
	s.NotEqual(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestFindOneFrom() {
	person, err := s.q.FindOneFrom(PersonTable, "id", 102)
	s.NoError(err)
	s.Equal(
		&Person{ID: 102, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", Email: pointer.ToString("elfrieda_abbott@example.org"), CreatedAt: personCreated},
		person,
	)

	project, err := s.q.FindOneFrom(ProjectTable, "id", "queen")
	s.NoError(err)
	s.Equal(&Project{ID: "queen", Name: "Thirsty Queen", Start: queenStart}, project)

	project, err = s.q.FindOneFrom(ProjectTable, "id", nil)
	s.Nil(project)
	s.Equal(reform.ErrNoRows, err)

	project, err = s.q.FindOneFrom(ProjectTable, "invalid_column", nil)
	s.Nil(project)
	s.Error(err)
	s.NotEqual(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestFindRows() {
	rows, err := s.q.FindRows(PersonTable, "name", "Elfrieda Abbott")
	s.NotNil(rows)
	s.NoError(err)
	defer rows.Close()

	var person Person
	err = s.q.NextRow(&person, rows)
	s.NoError(err)
	expected := Person{ID: 102, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", Email: pointer.ToString("elfrieda_abbott@example.org"), CreatedAt: personCreated}
	s.Equal(expected, person)

	err = s.q.NextRow(&person, rows)
	s.NoError(err)
	expected = Person{ID: 103, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", CreatedAt: personCreated}
	s.Equal(expected, person)

	err = s.q.NextRow(&person, rows)
	s.Equal(err, reform.ErrNoRows)
	s.Equal(expected, person) // expect old value

	rows, err = s.q.FindRows(ProjectTable, "id", nil)
	s.NotNil(rows)
	s.NoError(err)
	defer rows.Close()

	var project Project
	err = s.q.NextRow(&project, rows)
	s.Equal(reform.ErrNoRows, err)
	s.Equal(Project{}, project)

	rows, err = s.q.FindRows(ProjectTable, "invalid_column", nil)
	s.Nil(rows)
	s.Error(err)
	s.NotEqual(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestFindAllFrom() {
	structs, err := s.q.FindAllFrom(PersonTable, "name", "Elfrieda Abbott")
	s.NoError(err)
	s.Len(structs, 2)
	s.Equal([]reform.Struct{
		&Person{ID: 102, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", Email: pointer.ToString("elfrieda_abbott@example.org"), CreatedAt: personCreated},
		&Person{ID: 103, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", CreatedAt: personCreated},
	}, structs)

	structs, err = s.q.FindAllFrom(PersonTable, "id", 102, 103)
	s.NoError(err)
	s.Len(structs, 2)
	s.Equal([]reform.Struct{
		&Person{ID: 102, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", Email: pointer.ToString("elfrieda_abbott@example.org"), CreatedAt: personCreated},
		&Person{ID: 103, GroupID: pointer.ToInt32(65534), Name: "Elfrieda Abbott", CreatedAt: personCreated},
	}, structs)

	structs, err = s.q.FindAllFrom(ProjectTable, "id", nil)
	s.Nil(structs)
	s.NoError(err)

	structs, err = s.q.FindAllFrom(ProjectTable, "invalid_column", nil)
	s.Nil(structs)
	s.Error(err)
	s.NotEqual(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestFindByPrimaryKeyTo() {
	var person Person
	err := s.q.FindByPrimaryKeyTo(&person, 1)
	s.NoError(err)
	s.Equal(Person{ID: 1, GroupID: pointer.ToInt32(65534), Name: "Denis Mills", CreatedAt: goCreated}, person)

	var project Project
	err = s.q.FindByPrimaryKeyTo(&project, "baron")
	s.NoError(err)
	expected := Project{ID: "baron", Name: "Vicious Baron", Start: baronStart, End: &baronEnd}
	s.Equal(expected, project)

	err = s.q.FindByPrimaryKeyTo(&project, nil)
	s.Equal(expected, project) // expect old value
	s.Equal(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestFindByPrimaryKeyFrom() {
	person, err := s.q.FindByPrimaryKeyFrom(PersonTable, 1)
	s.NoError(err)
	s.Equal(&Person{ID: 1, GroupID: pointer.ToInt32(65534), Name: "Denis Mills", CreatedAt: goCreated}, person)

	project, err := s.q.FindByPrimaryKeyFrom(ProjectTable, "baron")
	s.NoError(err)
	s.Equal(&Project{ID: "baron", Name: "Vicious Baron", Start: baronStart, End: &baronEnd}, project)

	project, err = s.q.FindByPrimaryKeyFrom(ProjectTable, nil)
	s.Nil(project)
	s.Equal(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestReload() {
	person := Person{ID: 1}
	err := s.q.Reload(&person)
	s.NoError(err)
	s.Equal(Person{ID: 1, GroupID: pointer.ToInt32(65534), Name: "Denis Mills", CreatedAt: goCreated}, person)

	project := Project{ID: "baron"}
	err = s.q.Reload(&project)
	s.NoError(err)
	expected := Project{ID: "baron", Name: "Vicious Baron", Start: baronStart, End: &baronEnd}
	s.Equal(expected, project)

	project = Project{}
	err = s.q.Reload(&project)
	s.Equal(Project{}, project) // expect old value
	s.Equal(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestSelectsSchema() {
	if s.q.Dialect != postgresql.Dialect {
		s.T().Skip("only PostgreSQL supports schemas")
	}

	var legacyPerson LegacyPerson
	err := s.q.SelectOneTo(&legacyPerson, "WHERE id = "+s.q.Placeholder(1), 1001)
	s.NoError(err)
	s.Equal(LegacyPerson{ID: 1001, Name: pointer.ToString("Amelia Heathcote")}, legacyPerson)

	structs, err := s.q.FindAllFrom(LegacyPersonTable, "id", 1002, 1003)
	s.NoError(err)
	s.Len(structs, 2)
	s.Equal([]reform.Struct{
		&LegacyPerson{ID: 1002, Name: pointer.ToString("Anastacio Ledner")},
		&LegacyPerson{ID: 1003, Name: pointer.ToString("Dena Cummings")},
	}, structs)
}
