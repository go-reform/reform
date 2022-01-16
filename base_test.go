package reform_test

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql" //nolint:staticcheck
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
	"gopkg.in/reform.v1/internal/test"
	. "gopkg.in/reform.v1/internal/test/models"
)

// DB is a global connection pool shared by tests and examples.
//
// Deprecated: do not add new tests using it as using a global pool makes tests more brittle.
var DB *reform.DB

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		log.Print("Not setting DB in short mode")
	} else {
		DB = test.ConnectToTestDB()
	}

	os.Exit(m.Run())
}

// checkForeignKeys checks that foreign keys are still enforced for sqlite3.
func checkForeignKeys(t testing.TB, q *reform.Querier) {
	t.Helper()

	if q.Dialect != sqlite3.Dialect {
		return
	}

	var enabled bool
	err := q.QueryRow("PRAGMA foreign_keys").Scan(&enabled)
	require.NoError(t, err)
	require.True(t, enabled)
}

// withIdentityInsert executes an action with MS SQL IDENTITY_INSERT enabled for a table
func withIdentityInsert(t testing.TB, q *reform.Querier, table string, action func()) {
	t.Helper()

	if q.Dialect != mssql.Dialect && q.Dialect != sqlserver.Dialect { //nolint:staticcheck
		action()
		return
	}

	query := fmt.Sprintf("SET IDENTITY_INSERT %s %%s", q.QuoteIdentifier(table))

	_, err := q.Exec(fmt.Sprintf(query, "ON"))
	require.NoError(t, err)

	action()

	_, err = q.Exec(fmt.Sprintf(query, "OFF"))
	require.NoError(t, err)
}

func insertPersonWithID(t testing.TB, q *reform.Querier, str reform.Struct) error {
	t.Helper()

	var err error
	withIdentityInsert(t, q, "people", func() { err = q.Insert(str) })
	return err
}

// setupDB creates new database connection pool.
func setupDB(t testing.TB) *reform.DB {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	db := test.ConnectToTestDB()
	pl := reform.NewPrintfLogger(t.Logf)
	pl.LogTypes = true
	db.Logger = pl
	db.Querier = db.WithTag("test:%s", t.Name())

	checkForeignKeys(t, db.Querier)
	return db
}

// setupTX creates new database connection pool and starts a new transaction.
func setupTX(t testing.TB) (*reform.DB, *reform.TX) {
	t.Helper()

	db := setupDB(t)

	tx, err := db.Begin()
	require.NoError(t, err)
	return db, tx
}

// teardown closes database connection pool.
func teardown(t testing.TB, db *reform.DB) {
	t.Helper()

	err := db.DBInterface().(*sql.DB).Close()
	require.NoError(t, err)
}

// Deprecated: do not add new test to this suite, use Go subtests instead.
// TODO Remove.
type ReformSuite struct {
	suite.Suite
	tx *reform.TX
	q  *reform.Querier
}

func TestReformSuite(t *testing.T) {
	suite.Run(t, new(ReformSuite))
}

// SetupTest configures global connection pool and starts a new transaction.
func (s *ReformSuite) SetupTest() {
	if testing.Short() {
		s.T().Skip("skipping in short mode")
	}

	pl := reform.NewPrintfLogger(s.T().Logf)
	pl.LogTypes = true
	DB.Logger = pl
	DB.Querier = DB.WithTag("test:%s", s.T().Name())

	checkForeignKeys(s.T(), DB.Querier)

	tx, err := DB.Begin()
	s.Require().NoError(err)
	s.tx = tx
	s.q = tx.Querier
}

// TearDownTest rollbacks transaction created by SetupTest.
func (s *ReformSuite) TearDownTest() {
	if s.tx == nil {
		panic(s.T().Name() + ": tx is nil")
	}
	if s.q == nil {
		panic(s.T().Name() + ": q is nil")
	}

	checkForeignKeys(s.T(), s.q)
	s.Require().NoError(s.tx.Rollback())

	DB.Logger = nil
	DB.Querier = DB.WithTag("")
}

func (s *ReformSuite) RestartTransaction() {
	s.TearDownTest()
	s.SetupTest()
}

func (s *ReformSuite) TestStringer() {
	person, err := s.q.FindByPrimaryKeyFrom(PersonTable, 1)
	s.NoError(err)
	expected := "ID: 1 (int32), GroupID: 65534 (*int32), Name: `Denis Mills` (string), Email: <nil> (*string), CreatedAt: 2009-11-10 23:00:00 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)"
	s.Equal(expected, person.String())

	project, err := s.q.FindByPrimaryKeyFrom(ProjectTable, "baron")
	s.NoError(err)
	expected = "Name: `Vicious Baron` (string), ID: `baron` (string), Start: 2014-06-01 00:00:00 +0000 UTC (time.Time), End: 2016-02-21 00:00:00 +0000 UTC (*time.Time)"
	s.Equal(expected, project.String())
}

func (s *ReformSuite) TestNeverNil() {
	project := new(Project)

	for i, v := range project.Values() {
		if v == nil {
			s.Fail("Value is nil", "%s %#v", ProjectTable.Columns()[i], v)
		}
	}

	for i, v := range project.Pointers() {
		if v == nil {
			s.Fail("Pointer is nil", "%s %#v", ProjectTable.Columns()[i], v)
		}
	}

	v := project.PKValue()
	if v == nil {
		s.Fail("PKValue is nil")
	}

	v = project.PKPointer()
	if v == nil {
		s.Fail("PKPointer is nil")
	}
}

func (s *ReformSuite) TestHasPK() {
	person := new(Person)
	project := new(Project)
	s.False(person.HasPK())
	s.False(project.HasPK())

	person.ID = 1
	project.ID = "id"
	s.True(person.HasPK())
	s.True(project.HasPK())
}

func (s *ReformSuite) TestPlaceholders() {
	if s.q.Dialect != postgresql.Dialect {
		s.T().Skip("PostgreSQL-specific test")
	}

	s.Equal([]string{"$1", "$2", "$3", "$4", "$5"}, s.q.Placeholders(1, 5))
	s.Equal([]string{"$2", "$3", "$4", "$5", "$6"}, s.q.Placeholders(2, 5))
}

func (s *ReformSuite) TestTimezones() {
	t1 := time.Now().Round(0)
	t2 := t1.UTC()
	vlat, err := time.LoadLocation("Asia/Vladivostok")
	s.NoError(err)
	tVLAT := t1.In(vlat)
	hst, err := time.LoadLocation("US/Hawaii")
	s.NoError(err)
	tHST := t1.In(hst)

	{
		q := fmt.Sprintf(`INSERT INTO people (id, name, created_at) VALUES `+
			`(11, '11', %s), (12, '12', %s), (13, '13', %s), (14, '14', %s)`,
			s.q.Placeholder(1), s.q.Placeholder(2), s.q.Placeholder(3), s.q.Placeholder(4))

		withIdentityInsert(s.T(), s.q, "people", func() {
			_, err := s.q.Exec(q, t1, t2, tVLAT, tHST)
			s.NoError(err)
		})

		q = `SELECT created_at, created_at FROM people WHERE id IN (11, 12, 13, 14) ORDER BY id`
		rows, err := s.q.Query(q)
		s.NoError(err)

		for _, t := range []time.Time{t1, t2, tVLAT, tHST} {
			var createdS string
			var createdT time.Time
			rows.Next()
			err = rows.Scan(&createdS, &createdT)
			s.NoError(err)
			log.Printf("%s read from database as %q and %s", t, createdS, createdT)
		}

		s.NoError(rows.Err())
		s.NoError(rows.Close())
	}

	{
		q := fmt.Sprintf(`INSERT INTO projects (id, name, start) VALUES `+
			`('11', '11', %s), ('12', '12', %s), ('13', '13', %s), ('14', '14', %s)`,
			s.q.Placeholder(1), s.q.Placeholder(2), s.q.Placeholder(3), s.q.Placeholder(4))
		_, err := s.q.Exec(q, t1, t2, tVLAT, tHST)
		s.NoError(err)

		q = `SELECT start, start FROM projects WHERE id IN ('11', '12', '13', '14') ORDER BY id`
		rows, err := s.q.Query(q)
		s.NoError(err)

		for _, t := range []time.Time{t1, t2, tVLAT, tHST} {
			var startS string
			var startT time.Time
			rows.Next()
			err = rows.Scan(&startS, &startT)
			s.NoError(err)
			log.Printf("%s read from database as %q and %s", t, startS, startT)
		}

		s.NoError(rows.Err())
		s.NoError(rows.Close())
	}
}

// database/sql.(*Rows).Columns() is not currently used, but may be useful in the future.
// Test is in place to track drivers supporting it.
func (s *ReformSuite) TestColumns() {
	rows, err := s.q.SelectRows(PersonTable, "WHERE name = "+s.q.Placeholder(1)+" ORDER BY id", "Elfrieda Abbott")
	s.NoError(err)
	s.Require().NotNil(rows)
	s.NoError(rows.Err())

	columns, err := rows.Columns()
	s.NoError(err)
	s.Equal(PersonTable.Columns(), columns)
	s.NoError(rows.Close())
}

//nolint:staticcheck
func TestSetPK(t *testing.T) {
	t.Parallel()

	var person Person
	person.SetPK(int32(1))
	assert.EqualValues(t, 1, person.ID)
	person.SetPK(int64(2))
	assert.EqualValues(t, 2, person.ID)
	person.SetPK(Integer(3))
	assert.EqualValues(t, 3, person.ID)

	var project Project
	project.SetPK("baron")
	assert.EqualValues(t, "baron", project.ID)
	project.SetPK(1)
	assert.EqualValues(t, "baron", project.ID)

	var extra Extra
	extra.SetPK(int32(1))
	assert.EqualValues(t, 1, extra.ID)
	extra.SetPK(int64(2))
	assert.EqualValues(t, 2, extra.ID)
	extra.SetPK(Integer(3))
	assert.EqualValues(t, 3, extra.ID)
}
