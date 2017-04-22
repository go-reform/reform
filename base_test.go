package reform_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
	"gopkg.in/reform.v1/internal"
	. "gopkg.in/reform.v1/internal/test/models"
)

var (
	DB *reform.DB
)

func TestMain(m *testing.M) {
	DB = internal.ConnectToTestDB()
	os.Exit(m.Run())
}

// checkForeignKeys checks that foreign keys are still enforced for sqlite3.
func checkForeignKeys(t *testing.T, q *reform.Querier) {
	if q.Dialect != sqlite3.Dialect {
		return
	}

	var enabled bool
	err := q.QueryRow("PRAGMA foreign_keys").Scan(&enabled)
	require.NoError(t, err)
	require.True(t, enabled)
}

// setIdentityInsert allows or disallows insertions of rows with set primary keys for MS SQL.
func setIdentityInsert(t *testing.T, q *reform.Querier, table string, allow bool) {
	if q.Dialect != mssql.Dialect && q.Dialect != sqlserver.Dialect {
		return
	}

	allowString := "OFF"
	if allow {
		allowString = "ON"
	}
	sql := fmt.Sprintf("SET IDENTITY_INSERT %s %s", q.QuoteIdentifier(table), allowString)
	_, err := q.Exec(sql)
	require.NoError(t, err)
}

type ReformSuite struct {
	suite.Suite
	tx *reform.TX
	q  *reform.Querier
}

func TestReformSuite(t *testing.T) {
	suite.Run(t, new(ReformSuite))
}

func (s *ReformSuite) SetupTest() {
	pl := reform.NewPrintfLogger(s.T().Logf)
	pl.LogTypes = true
	DB.Logger = pl

	var err error
	s.tx, err = DB.Begin()
	s.Require().NoError(err)

	s.q = s.tx.WithTag("test")

	setIdentityInsert(s.T(), s.q, "people", false)
}

func (s *ReformSuite) TearDownTest() {
	// some transactional tests rollback and nilify q
	if s.q != nil {
		checkForeignKeys(s.T(), s.q)

		err := s.tx.Rollback()
		s.Require().NoError(err)
	}

	checkForeignKeys(s.T(), DB.Querier)
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
	setIdentityInsert(s.T(), s.q, "people", true)

	t1 := time.Now()
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
		_, err := s.q.Exec(q, t1, t2, tVLAT, tHST)
		s.NoError(err)

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

		err = rows.Close()
		s.NoError(err)
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
		defer rows.Close()

		for _, t := range []time.Time{t1, t2, tVLAT, tHST} {
			var startS string
			var startT time.Time
			rows.Next()
			err = rows.Scan(&startS, &startT)
			s.NoError(err)
			log.Printf("%s read from database as %q and %s", t, startS, startT)
		}

		err = rows.Close()
		s.NoError(err)
	}
}

// database/sql.(*Rows).Columns() is not currently used, but may be useful in the future.
// Test is in place to track drivers supporting it.
func (s *ReformSuite) TestColumns() {
	rows, err := s.q.SelectRows(PersonTable, "WHERE name = "+s.q.Placeholder(1)+" ORDER BY id", "Elfrieda Abbott")
	s.NoError(err)
	s.Require().NotNil(rows)
	defer rows.Close()

	columns, err := rows.Columns()
	s.NoError(err)
	s.Equal(PersonTable.Columns(), columns)
}
