package reform_test

import (
	"database/sql"
	"errors"
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
	_ "github.com/ziutek/mymysql/godrv"

	"github.com/AlekSi/pointer"
	"github.com/enodata/faker"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/internal/test/models"
)

var (
	DB *reform.DB
)

func TestMain(m *testing.M) {
	driver := os.Getenv("REFORM_TEST_DRIVER")
	source := os.Getenv("REFORM_TEST_SOURCE")
	log.Printf("driver = %q, source = %q", driver, source)
	if driver == "" || source == "" {
		log.Fatal("no driver or source, set REFORM_TEST_DRIVER and REFORM_TEST_SOURCE")
	}

	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(-1)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	log.Printf("time.Now()       = %s", now)
	log.Printf("time.Now().UTC() = %s", now.UTC())

	var dialect reform.Dialect
	switch driver {
	case "mysql", "mymysql":
		dialect = mysql.Dialect

		var tz string
		err = db.QueryRow("SHOW VARIABLES LIKE 'time_zone'").Scan(&tz, &tz)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("MySQL time_zone = %q", tz)

	case "postgres", "pgx":
		dialect = postgresql.Dialect

		var tz string
		err = db.QueryRow("SHOW TimeZone").Scan(&tz)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("PostgreSQL TimeZone = %q", tz)

	case "sqlite3":
		dialect = sqlite3.Dialect

		_, err = db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			log.Fatal(err)
		}

	case "mssql":
		dialect = mssql.Dialect

	default:
		log.Fatal("reform: no dialect for driver " + driver)
	}

	DB = reform.NewDB(db, dialect, nil)

	os.Exit(m.Run())
}

func setIdentityInsert(t *testing.T, tx *reform.TX, table string, allow bool) {
	if tx.Dialect != mssql.Dialect {
		return
	}

	allowString := "OFF"
	if allow {
		allowString = "ON"
	}
	sql := fmt.Sprintf("SET IDENTITY_INSERT %s %s", tx.QuoteIdentifier(table), allowString)
	_, err := tx.Exec(sql)
	require.NoError(t, err)
}

type ReformSuite struct {
	suite.Suite
	q *reform.TX
}

func TestReformSuite(t *testing.T) {
	suite.Run(t, new(ReformSuite))
}

func (s *ReformSuite) SetupTest() {
	pl := reform.NewPrintfLogger(s.T().Logf)
	pl.LogTypes = true
	DB.Logger = pl

	var err error
	s.q, err = DB.Begin()
	s.Require().NoError(err)

	setIdentityInsert(s.T(), s.q, "people", false)
}

func (s *ReformSuite) TearDownTest() {
	// some transactional tests rollback and nilify transaction
	if s.q != nil {
		err := s.q.Rollback()
		s.Require().NoError(err)
	}
}

func (s *ReformSuite) RestartTransaction() {
	s.TearDownTest()
	s.SetupTest()
}

func (s *ReformSuite) TestStringer() {
	person, err := s.q.FindByPrimaryKeyFrom(models.PersonTable, 1)
	s.NoError(err)
	expected := "ID: 1 (int32), GroupID: 65534 (*int32), Name: `Denis Mills` (string), Email: <nil> (*string), CreatedAt: 2009-11-10 23:00:00 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)"
	s.Equal(expected, person.String())

	project, err := s.q.FindByPrimaryKeyFrom(models.ProjectTable, "baron")
	s.NoError(err)
	expected = "Name: `Vicious Baron` (string), ID: `baron` (string), Start: 2014-06-01 00:00:00 +0000 UTC (time.Time), End: 2016-02-21 00:00:00 +0000 UTC (*time.Time)"
	s.Equal(expected, project.String())
}

func (s *ReformSuite) TestNeverNil() {
	project := new(models.Project)

	for i, v := range project.Values() {
		if v == nil {
			s.Fail("Value is nil", "%s %#v", models.ProjectTable.Columns()[i], v)
		}
	}

	for i, v := range project.Pointers() {
		if v == nil {
			s.Fail("Pointer is nil", "%s %#v", models.ProjectTable.Columns()[i], v)
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
	person := new(models.Person)
	project := new(models.Project)
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

func (s *ReformSuite) TestInTransaction() {
	setIdentityInsert(s.T(), s.q, "people", true)

	err := s.q.Rollback()
	s.Require().NoError(err)
	s.q = nil

	person := &models.Person{ID: 42, Email: pointer.ToString(faker.Internet().Email())}

	err = DB.InTransaction(func(tx *reform.TX) error {
		err := tx.Insert(person)
		s.NoError(err)
		return errors.New("epic error")
	})
	s.EqualError(err, "epic error")

	s.Panics(func() {
		err = DB.InTransaction(func(tx *reform.TX) error {
			err := tx.Insert(person)
			s.NoError(err)
			panic("epic panic!")
		})
	})

	err = DB.InTransaction(func(tx *reform.TX) error {
		err := tx.Insert(person)
		s.NoError(err)
		return nil
	})
	s.NoError(err)

	err = DB.Insert(person)
	s.Error(err)

	err = DB.Delete(person)
	s.NoError(err)
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
