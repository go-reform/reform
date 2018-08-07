package reform_test

import (
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
	. "gopkg.in/reform.v1/internal/test/models"
)

func (s *ReformSuite) TestQualifiedView() {
	switch s.q.Dialect {
	case postgresql.Dialect:
		s.Equal(`"people"`, s.q.QualifiedView(PersonTable))
		s.Equal(`"people"`, s.q.WithQualifiedViewName("ignored").QualifiedView(PersonTable))
		s.Equal(`"legacy"."people"`, s.q.QualifiedView(LegacyPersonTable))
		s.Equal(`"legacy"."people"`, s.q.WithQualifiedViewName("ignored").QualifiedView(LegacyPersonTable))

	case mysql.Dialect:
		s.Equal("`people`", s.q.QualifiedView(PersonTable))
		s.Equal("`people`", s.q.WithQualifiedViewName("ignored").QualifiedView(PersonTable))

	case sqlite3.Dialect:
		s.Equal(`"people"`, s.q.QualifiedView(PersonTable))
		s.Equal(`"people"`, s.q.WithQualifiedViewName("ignored").QualifiedView(PersonTable))

	case mssql.Dialect, sqlserver.Dialect:
		s.Equal(`[people]`, s.q.QualifiedView(PersonTable))
		s.Equal(`[people]`, s.q.WithQualifiedViewName("ignored").QualifiedView(PersonTable))

	default:
		s.Fail("Unhandled dialect", s.q.Dialect.String())
	}
}
