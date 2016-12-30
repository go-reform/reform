package reform_test

import (
	"errors"

	"github.com/AlekSi/pointer"
	"github.com/enodata/faker"

	reform "gopkg.in/reform.v1"
	. "gopkg.in/reform.v1/internal/test/models"
)

func (s *ReformSuite) TestInSavepoint() {
	setIdentityInsert(s.T(), s.q, "people", true)

	person := &Person{ID: 42, Email: pointer.ToString(faker.Internet().Email())}

	// error in closure
	err := s.q.InSavepoint(func() error {
		s.NoError(s.q.Insert(person))
		return errors.New("epic error")
	})
	s.EqualError(err, "epic error")
	s.Equal(s.q.Reload(person), reform.ErrNoRows)

	// panic in closure
	s.Panics(func() {
		err = s.q.InSavepoint(func() error {
			s.NoError(s.q.Insert(person))
			panic("epic panic!")
		})
	})
	s.Equal(s.q.Reload(person), reform.ErrNoRows)

	// duplicate PK in closure
	err = s.q.InSavepoint(func() error {
		s.NoError(s.q.Insert(person))
		err := s.q.Insert(person)
		s.Error(err)
		return err
	})
	s.Error(err)
	s.Equal(s.q.Reload(person), reform.ErrNoRows)

	// no error
	err = s.q.InSavepoint(func() error {
		s.NoError(s.q.Insert(person))
		return nil
	})
	s.NoError(err)
	s.NoError(s.q.Reload(person))
	s.NoError(s.q.Delete(person))
}
