package reform_test

import (
	"errors"

	"github.com/AlekSi/pointer"
	"github.com/enodata/faker"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	. "gopkg.in/reform.v1/internal/test/models"
)

func (s *ReformSuite) TestBeginCommit() {
	s.Require().NoError(s.tx.Rollback())
	s.q = nil

	setIdentityInsert(s.T(), DB.Querier, "people", true)

	person := &Person{ID: 42, Email: pointer.ToString(faker.Internet().Email())}

	tx, err := DB.Begin()
	s.Require().NoError(err)
	s.NoError(tx.Insert(person))
	s.NoError(tx.Commit())
	s.Equal(tx.Commit(), reform.ErrTxDone)
	s.Equal(tx.Rollback(), reform.ErrTxDone)
	s.NoError(DB.Reload(person))
	s.NoError(DB.Delete(person))
}

func (s *ReformSuite) TestBeginRollback() {
	s.Require().NoError(s.tx.Rollback())
	s.q = nil

	setIdentityInsert(s.T(), DB.Querier, "people", true)

	person := &Person{ID: 42, Email: pointer.ToString(faker.Internet().Email())}

	tx, err := DB.Begin()
	s.Require().NoError(err)
	s.NoError(tx.Insert(person))
	s.NoError(tx.Rollback())
	s.Equal(tx.Commit(), reform.ErrTxDone)
	s.Equal(tx.Rollback(), reform.ErrTxDone)
	s.Equal(DB.Reload(person), reform.ErrNoRows)
}

// This behavior is checked for documentation purposes only. reform does not rely on it.
func (s *ReformSuite) TestErrorInTransaction() {
	if s.q.Dialect == postgresql.Dialect {
		s.T().Skip(s.q.Dialect.String() + " works differently, see TestAbortedTransaction")
	}

	s.Require().NoError(s.tx.Rollback())
	s.q = nil

	setIdentityInsert(s.T(), DB.Querier, "people", true)

	person1 := &Person{ID: 42, Email: pointer.ToString(faker.Internet().Email())}
	person2 := &Person{ID: 43, Email: pointer.ToString(faker.Internet().Email())}

	// commit works
	tx, err := DB.Begin()
	s.Require().NoError(err)
	s.NoError(tx.Insert(person1))
	s.Error(tx.Insert(person1))   // duplicate PK
	s.NoError(tx.Insert(person2)) // INSERT works
	s.NoError(tx.Commit())
	s.Equal(tx.Commit(), reform.ErrTxDone)
	s.Equal(tx.Rollback(), reform.ErrTxDone)
	s.NoError(DB.Reload(person1))
	s.NoError(DB.Reload(person2))
	s.NoError(DB.Delete(person1))
	s.NoError(DB.Delete(person2))

	// rollback works
	tx, err = DB.Begin()
	s.Require().NoError(err)
	s.NoError(tx.Insert(person1))
	s.Error(tx.Insert(person1))   // duplicate PK
	s.NoError(tx.Insert(person2)) // INSERT works
	s.NoError(tx.Rollback())
	s.Equal(tx.Commit(), reform.ErrTxDone)
	s.Equal(tx.Rollback(), reform.ErrTxDone)
	s.EqualError(DB.Reload(person1), reform.ErrNoRows.Error())
	s.EqualError(DB.Reload(person2), reform.ErrNoRows.Error())
}

// This behavior is checked for documentation purposes only. reform does not rely on it.
// http://postgresql.nabble.com/Current-transaction-is-aborted-commands-ignored-until-end-of-transaction-block-td5109252.html
func (s *ReformSuite) TestAbortedTransaction() {
	if s.q.Dialect != postgresql.Dialect {
		s.T().Skip(s.q.Dialect.String() + " works differently, see TestErrorInTransaction")
	}

	s.Require().NoError(s.tx.Rollback())
	s.q = nil

	setIdentityInsert(s.T(), DB.Querier, "people", true)

	person1 := &Person{ID: 42, Email: pointer.ToString(faker.Internet().Email())}
	person2 := &Person{ID: 43, Email: pointer.ToString(faker.Internet().Email())}

	// commit fails
	tx, err := DB.Begin()
	s.Require().NoError(err)
	s.NoError(tx.Insert(person1))
	s.EqualError(tx.Insert(person1), `pq: duplicate key value violates unique constraint "people_pkey"`)
	s.EqualError(tx.Insert(person2), `pq: current transaction is aborted, commands ignored until end of transaction block`)
	s.EqualError(tx.Commit(), `pq: Could not complete operation in a failed transaction`)
	s.Equal(tx.Rollback(), reform.ErrTxDone)
	s.EqualError(DB.Reload(person1), reform.ErrNoRows.Error())
	s.EqualError(DB.Reload(person2), reform.ErrNoRows.Error())

	// rollback works
	tx, err = DB.Begin()
	s.Require().NoError(err)
	s.NoError(tx.Insert(person1))
	s.EqualError(tx.Insert(person1), `pq: duplicate key value violates unique constraint "people_pkey"`)
	s.EqualError(tx.Insert(person2), `pq: current transaction is aborted, commands ignored until end of transaction block`)
	s.NoError(tx.Rollback())
	s.Equal(tx.Commit(), reform.ErrTxDone)
	s.Equal(tx.Rollback(), reform.ErrTxDone)
	s.EqualError(DB.Reload(person1), reform.ErrNoRows.Error())
	s.EqualError(DB.Reload(person2), reform.ErrNoRows.Error())
}

func (s *ReformSuite) TestInTransaction() {
	s.Require().NoError(s.tx.Rollback())
	s.q = nil

	setIdentityInsert(s.T(), DB.Querier, "people", true)

	person := &Person{ID: 42, Email: pointer.ToString(faker.Internet().Email())}

	// error in closure
	err := DB.InTransaction(func(tx *reform.TX) error {
		s.NoError(tx.Insert(person))
		return errors.New("epic error")
	})
	s.EqualError(err, "epic error")
	s.Equal(DB.Reload(person), reform.ErrNoRows)

	// panic in closure
	s.Panics(func() {
		err = DB.InTransaction(func(tx *reform.TX) error {
			s.NoError(tx.Insert(person))
			panic("epic panic!")
		})
	})
	s.Equal(DB.Reload(person), reform.ErrNoRows)

	// duplicate PK in closure
	err = DB.InTransaction(func(tx *reform.TX) error {
		s.NoError(tx.Insert(person))
		err := tx.Insert(person)
		s.Error(err)
		return err
	})
	s.Error(err)
	s.Equal(DB.Reload(person), reform.ErrNoRows)

	// no error
	err = DB.InTransaction(func(tx *reform.TX) error {
		s.NoError(tx.Insert(person))
		return nil
	})
	s.NoError(err)
	s.NoError(DB.Reload(person))
	s.NoError(DB.Delete(person))
}
