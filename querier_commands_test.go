package reform_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/enodata/faker"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	. "gopkg.in/reform.v1/internal/test/models"
)

func (s *ReformSuite) TestInsert() {
	newEmail := faker.Internet().Email()
	person := &Person{Email: &newEmail}
	err := s.q.Insert(person)
	s.NoError(err)
	s.NotEqual(int32(0), person.ID)
	s.Equal("", person.Name)
	s.Equal(&newEmail, person.Email)
	s.WithinDuration(time.Now(), person.CreatedAt, 2*time.Second)
	s.Nil(person.UpdatedAt)

	person2, err := s.q.FindByPrimaryKeyFrom(PersonTable, person.ID)
	s.NoError(err)
	s.Equal(person, person2)

	err = s.q.Insert(person)
	s.Error(err)
}

func (s *ReformSuite) TestInsertWithValues() {
	t := time.Now()
	newEmail := faker.Internet().Email()
	person := &Person{Email: &newEmail, CreatedAt: t, UpdatedAt: &t}
	err := s.q.Insert(person)
	s.NoError(err)
	s.NotEqual(int32(0), person.ID)
	s.Equal("", person.Name)
	s.Equal(&newEmail, person.Email)
	s.WithinDuration(t, person.CreatedAt, 2*time.Second)
	s.WithinDuration(t, *person.UpdatedAt, 2*time.Second)

	person2, err := s.q.FindByPrimaryKeyFrom(PersonTable, person.ID)
	s.NoError(err)
	s.Equal(person, person2)

	err = s.q.Insert(person)
	s.Error(err)
}

func (s *ReformSuite) TestInsertWithPrimaryKey() {
	setIdentityInsert(s.T(), s.q, "people", true)

	newEmail := faker.Internet().Email()
	person := &Person{ID: 50, Email: &newEmail}
	err := s.q.Insert(person)
	s.NoError(err)
	s.Equal(int32(50), person.ID)
	s.Equal("", person.Name)
	s.Equal(&newEmail, person.Email)
	s.WithinDuration(time.Now(), person.CreatedAt, 2*time.Second)
	s.Nil(person.UpdatedAt)

	person2, err := s.q.FindByPrimaryKeyFrom(PersonTable, person.ID)
	s.NoError(err)
	s.Equal(person, person2)

	err = s.q.Insert(person)
	s.Error(err)
}

func (s *ReformSuite) TestInsertWithStringPrimaryKey() {
	start := time.Now().UTC().Truncate(24 * time.Hour)
	project := &Project{ID: "new", Start: start, End: pointer.ToTime(start.AddDate(0, 0, 1))}
	err := s.q.Insert(project)
	s.NoError(err)
	s.Equal("new", project.ID)

	project2, err := s.q.FindByPrimaryKeyFrom(ProjectTable, project.ID)
	s.NoError(err)
	s.Equal(project, project2)

	err = s.q.Insert(project)
	s.Error(err)
}

func (s *ReformSuite) TestInsertIntoView() {
	pp := &PersonProject{PersonID: 1, ProjectID: "baron"}
	err := s.q.Insert(pp)
	s.NoError(err)

	err = s.q.Insert(pp)
	s.Error(err)

	s.RestartTransaction()

	pp = &PersonProject{PersonID: 1, ProjectID: "no_such_project"}
	err = s.q.Insert(pp)
	s.Error(err)
}

func (s *ReformSuite) TestInsertColumns() {
	t := time.Now()
	newEmail := faker.Internet().Email()
	person := &Person{Email: &newEmail, CreatedAt: t, UpdatedAt: &t}
	columns := []string{"name", "email", "created_at", "updated_at"}
	err := s.q.InsertColumns(person, columns...)
	s.NoError(err)
	s.NotEqual(int32(0), person.ID)
	s.Equal("", person.Name)
	s.Equal((*int32)(nil), person.GroupID)
	s.Equal(&newEmail, person.Email)
	s.WithinDuration(t, person.CreatedAt, 2*time.Second)
	s.WithinDuration(t, *person.UpdatedAt, 2*time.Second)
	s.Equal([]string{"name", "email", "created_at", "updated_at"}, columns, "should not be changed")

	person2, err := s.q.FindByPrimaryKeyFrom(PersonTable, person.ID)
	s.NoError(err)
	person.GroupID = pointer.ToInt32(65534)
	s.Equal(person, person2)

	err = s.q.Insert(person)
	s.Error(err)
}

func (s *ReformSuite) TestInsertColumnsIntoView() {
	pp := &PersonProject{PersonID: 1, ProjectID: "baron"}
	err := s.q.InsertColumns(pp, "person_id", "project_id")
	s.NoError(err)

	err = s.q.InsertColumns(pp, "person_id", "project_id")
	s.Error(err)

	s.RestartTransaction()

	pp = &PersonProject{PersonID: 1, ProjectID: "no_such_project"}
	err = s.q.InsertColumns(pp, "person_id", "project_id")
	s.Error(err)
}

func (s *ReformSuite) TestInsertMulti() {
	newEmail := faker.Internet().Email()
	newName := faker.Name().Name()
	person1, person2 := &Person{Email: &newEmail}, &Person{Name: newName}
	err := s.q.InsertMulti(person1, person2)
	s.NoError(err)

	s.Equal(int32(0), person1.ID)
	s.Equal("", person1.Name)
	s.Equal(&newEmail, person1.Email)
	s.WithinDuration(time.Now(), person1.CreatedAt, 2*time.Second)
	s.Nil(person1.UpdatedAt)

	s.Equal(int32(0), person2.ID)
	s.Equal(newName, person2.Name)
	s.Nil(person2.Email)
	s.WithinDuration(time.Now(), person2.CreatedAt, 2*time.Second)
	s.Nil(person2.UpdatedAt)
}

func (s *ReformSuite) TestInsertMultiWithPrimaryKeys() {
	setIdentityInsert(s.T(), s.q, "people", true)

	newEmail := faker.Internet().Email()
	newName := faker.Name().Name()
	person1, person2 := &Person{ID: 50, Email: &newEmail}, &Person{ID: 51, Name: newName}
	err := s.q.InsertMulti(person1, person2)
	s.NoError(err)

	s.Equal(int32(50), person1.ID)
	s.Equal("", person1.Name)
	s.Equal(&newEmail, person1.Email)
	s.WithinDuration(time.Now(), person1.CreatedAt, 2*time.Second)
	s.Nil(person1.UpdatedAt)

	s.Equal(int32(51), person2.ID)
	s.Equal(newName, person2.Name)
	s.Nil(person2.Email)
	s.WithinDuration(time.Now(), person2.CreatedAt, 2*time.Second)
	s.Nil(person2.UpdatedAt)

	person, err := s.q.FindByPrimaryKeyFrom(PersonTable, person1.ID)
	s.NoError(err)
	s.Equal(person1, person)

	person, err = s.q.FindByPrimaryKeyFrom(PersonTable, person2.ID)
	s.NoError(err)
	s.Equal(person2, person)
}

func (s *ReformSuite) TestInsertMultiMixes() {
	err := s.q.InsertMulti()
	s.NoError(err)

	err = s.q.InsertMulti(&Person{}, &Project{})
	s.Error(err)

	err = s.q.InsertMulti(&Person{ID: 1}, &Person{})
	s.Error(err)
}

func (s *ReformSuite) TestInsertIDOnly() {
	id := &IDOnly{}
	err := s.q.Insert(id)
	s.NoError(err)
	s.Equal(int32(1), id.ID)
}

func (s *ReformSuite) TestUpdate() {
	var person Person
	err := s.q.Update(&person)
	s.Equal(reform.ErrNoPK, err)

	person.ID = 99
	err = s.q.Update(&person)
	s.Equal(reform.ErrNoRows, err)

	err = s.q.FindByPrimaryKeyTo(&person, 102)
	s.NoError(err)

	person.Email = pointer.ToString(faker.Internet().Email())
	err = s.q.Update(&person)
	s.NoError(err)
	s.Equal(personCreated, person.CreatedAt)
	s.Require().NotNil(person.UpdatedAt)
	s.WithinDuration(time.Now(), *person.UpdatedAt, 2*time.Second)

	person2, err := s.q.FindByPrimaryKeyFrom(PersonTable, person.ID)
	s.NoError(err)
	s.Equal(&person, person2)
}

func (s *ReformSuite) TestUpdateOverwrite() {
	newEmail := faker.Internet().Email()
	person := Person{ID: 102, Email: pointer.ToString(newEmail)}
	err := s.q.Update(&person)
	s.NoError(err)

	var person2 Person
	err = s.q.FindByPrimaryKeyTo(&person2, person.ID)
	s.NoError(err)
	s.Equal("", person2.Name)
	s.Equal(&newEmail, person2.Email)
	s.WithinDuration(time.Now(), person2.CreatedAt, 2*time.Second)
	s.Require().NotNil(person2.UpdatedAt)
	s.WithinDuration(time.Now(), *person2.UpdatedAt, 2*time.Second)
}

func (s *ReformSuite) TestUpdateColumns() {
	newName := faker.Name().Name()
	newEmail := faker.Internet().Email()

	for p, columns := range map[*Person][]string{
		&Person{Name: "Elfrieda Abbott", Email: &newEmail}:                             {"email", "updated_at"},
		&Person{Name: newName, Email: pointer.ToString("elfrieda_abbott@example.org")}: {"name", "name", "updated_at"},
		&Person{Name: newName, Email: &newEmail}:                                       {"name", "email", "updated_at"},
	} {
		var person Person
		err := s.q.FindByPrimaryKeyTo(&person, 102)
		s.NoError(err)

		columnsCopy := make([]string, len(columns))
		copy(columnsCopy, columns)

		person.Name = p.Name
		person.Email = p.Email
		err = s.q.UpdateColumns(&person, columns...)
		s.NoError(err)
		s.Equal(personCreated, person.CreatedAt)
		s.Require().NotNil(person.UpdatedAt)
		s.WithinDuration(time.Now(), *person.UpdatedAt, 2*time.Second)
		s.Equal(columnsCopy, columns, "should not be changed")

		person2, err := s.q.FindByPrimaryKeyFrom(PersonTable, person.ID)
		s.NoError(err)
		s.Equal(&person, person2)

		s.RestartTransaction()
	}

	person := &Person{ID: 102, Name: newName, Email: &newEmail, CreatedAt: personCreated}
	for e, columns := range map[error][]string{
		errors.New("reform: unexpected columns: [foo]"):     {"foo"},
		errors.New("reform: will not update PK column: id"): {"id"},
		errors.New("reform: nothing to update"):             {},
	} {
		err := s.q.UpdateColumns(person, columns...)
		s.Error(err)
		s.Equal(e, err)
	}
}

func (s *ReformSuite) TestUpdateView() {
	newName := faker.Name().Name()
	newEmail := faker.Internet().Email()

	for p, columns := range map[*Person][]string{
		&Person{Name: "Elfrieda Abbott", Email: &newEmail}:                             {"email", "updated_at"},
		&Person{Name: newName, Email: pointer.ToString("elfrieda_abbott@example.org")}: {"name", "name", "updated_at"},
		&Person{Name: newName, Email: &newEmail}:                                       {"name", "email", "updated_at"},
	} {
		var person Person
		err := s.q.FindByPrimaryKeyTo(&person, 102)
		s.NoError(err)

		columnsCopy := make([]string, len(columns))
		copy(columnsCopy, columns)

		person.Name = p.Name
		person.Email = p.Email
		ra, err := s.q.UpdateView(&person, columns, fmt.Sprintf("WHERE id = %d", person.ID))
		s.NoError(err)
		s.Equal(uint(1), ra)
		s.Equal(personCreated, person.CreatedAt)
		s.Require().NotNil(person.UpdatedAt)
		s.WithinDuration(time.Now(), *person.UpdatedAt, 2*time.Second)
		s.Equal(columnsCopy, columns, "should not be changed")

		person2, err := s.q.FindByPrimaryKeyFrom(PersonTable, person.ID)
		s.NoError(err)
		s.Equal(&person, person2)

		s.RestartTransaction()
	}

	person := &Person{ID: 102, Name: newName, Email: &newEmail, CreatedAt: personCreated}
	for e, columns := range map[error][]string{
		errors.New("reform: unexpected columns: [foo]"):     {"foo"},
		errors.New("reform: will not update PK column: id"): {"id"},
		errors.New("reform: nothing to update"):             {},
	} {
		ra, err := s.q.UpdateView(person, columns, "")
		s.Error(err)
		s.Zero(ra)
		s.Equal(e, err)
	}
}

func (s *ReformSuite) TestSave() {
	newName := faker.Name().Name()
	person := &Person{Name: newName}
	err := s.q.Save(person)
	s.NoError(err)

	person2, err := s.q.FindByPrimaryKeyFrom(PersonTable, person.ID)
	s.Require().NoError(err)
	s.Equal(newName, person2.(*Person).Name)
	s.Nil(person2.(*Person).Email)
	s.Equal(person, person2)

	newEmail := faker.Internet().Email()
	person.Email = &newEmail
	err = s.q.Save(person)
	s.NoError(err)

	person2, err = s.q.FindByPrimaryKeyFrom(PersonTable, person.ID)
	s.Require().NoError(err)
	s.Equal(newName, person2.(*Person).Name)
	s.Equal(&newEmail, person2.(*Person).Email)
	s.Equal(person, person2)
}

func (s *ReformSuite) TestSaveWithPrimaryKey() {
	setIdentityInsert(s.T(), s.q, "people", true)

	newName := faker.Name().Name()
	person := &Person{ID: 99, Name: newName}
	err := s.q.Save(person)
	s.NoError(err)

	// that should cause no-op UPDATE, see https://github.com/go-reform/reform/issues/131
	err = s.q.Save(person)
	s.NoError(err)
}

func (s *ReformSuite) TestDelete() {
	person := &Person{ID: 1}
	err := s.q.Delete(person)
	s.NoError(err)
	err = s.q.Reload(person)
	s.Equal(reform.ErrNoRows, err)

	project := &Project{ID: "baron"}
	err = s.q.Delete(project)
	s.NoError(err)
	err = s.q.Reload(project)
	s.Equal(reform.ErrNoRows, err)

	project = &Project{}
	err = s.q.Delete(project)
	s.Equal(reform.ErrNoPK, err)

	project = &Project{ID: "no_such_project"}
	err = s.q.Delete(project)
	s.Equal(reform.ErrNoRows, err)
}

func (s *ReformSuite) TestDeleteFrom() {
	ra, err := s.q.DeleteFrom(PersonTable, "WHERE email IS NULL")
	s.NoError(err)
	s.Equal(uint(3), ra)

	ra, err = s.q.DeleteFrom(PersonTable, "WHERE email IS NULL")
	s.NoError(err)
	s.Equal(uint(0), ra)

	// -1 second for SQLite3, otherwise it also deletes queen itself ¯\_(ツ)_/¯
	ra, err = s.q.DeleteFrom(ProjectTable, "WHERE start < "+s.q.Placeholder(1), queenStart.Add(-time.Second))
	s.NoError(err)
	s.Equal(uint(3), ra)

	ra, err = s.q.DeleteFrom(ProjectTable, "")
	s.NoError(err)
	s.Equal(uint(2), ra)

	ra, err = s.q.DeleteFrom(ProjectTable, "WHERE invalid_tail")
	s.Error(err)
	s.Equal(uint(0), ra)
}

func (s *ReformSuite) TestCommandsSchema() {
	if s.q.Dialect != postgresql.Dialect {
		s.T().Skip("only PostgreSQL supports schemas")
	}

	legacyPerson := &LegacyPerson{Name: pointer.ToString(faker.Name().Name())}
	err := s.q.Save(legacyPerson)
	s.NoError(err)
	err = s.q.Save(legacyPerson)
	s.NoError(err)
	err = s.q.Delete(legacyPerson)
	s.NoError(err)
}
