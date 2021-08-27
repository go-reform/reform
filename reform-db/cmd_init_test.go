package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/parse"
)

func (s *ReformDBSuite) TestInit() {
	good, err := parse.File("../internal/test/models/good.go")
	s.Require().NoError(err)
	s.Require().Len(good, 7)

	people := good[0]
	projects := good[1]
	personProject := good[2]
	idOnly := good[3]
	constraints := good[4]
	compositePK := good[5]

	// patch difference we don't handle
	people.Type = strings.ReplaceAll(people.Type, "Person", "People")
	projects.Type = strings.ReplaceAll(projects.Type, "Project", "Projects")
	if s.db.Dialect == sqlite3.Dialect {
		people.Fields[0].Type = strings.ReplaceAll(people.Fields[0].Type, "int32", "int64")
		people.Fields[1].Type = strings.ReplaceAll(people.Fields[1].Type, "int32", "int64")
		personProject.Fields[0].Type = strings.ReplaceAll(personProject.Fields[0].Type, "int32", "int64")
		idOnly.Fields[0].Type = strings.ReplaceAll(idOnly.Fields[0].Type, "int32", "int64")
		constraints.Fields[0].Type = strings.ReplaceAll(constraints.Fields[0].Type, "int32", "int64")
		compositePK.Fields[0].Type = strings.ReplaceAll(compositePK.Fields[0].Type, "int32", "int64")
	}

	dir, err := ioutil.TempDir("", "ReformDBTestInit")
	s.Require().NoError(err)
	s.T().Log(dir)

	cmdInit(s.db, dir)

	fis, err := ioutil.ReadDir(dir)
	s.Require().NoError(err)
	s.Require().Len(fis, 6)

	ff := filepath.Join(dir, "people.go")
	actual, err := parse.File(ff)
	s.Require().NoError(err)
	s.Require().Len(actual, 1)
	s.Require().Equal(people, actual[0])

	ff = filepath.Join(dir, "projects.go")
	actual, err = parse.File(ff)
	s.Require().NoError(err)
	s.Require().Len(actual, 1)
	s.Require().Equal(projects, actual[0])

	ff = filepath.Join(dir, "person_project.go")
	actual, err = parse.File(ff)
	s.Require().NoError(err)
	s.Require().Len(actual, 1)
	s.Require().Equal(personProject, actual[0])

	ff = filepath.Join(dir, "id_only.go")
	actual, err = parse.File(ff)
	s.Require().NoError(err)
	s.Require().Len(actual, 1)
	s.Require().Equal(idOnly, actual[0])

	ff = filepath.Join(dir, "constraints.go")
	actual, err = parse.File(ff)
	s.Require().NoError(err)
	s.Require().Len(actual, 1)
	s.Require().Equal(constraints, actual[0])

	ff = filepath.Join(dir, "composite_pk.go")
	actual, err = parse.File(ff)
	s.Require().NoError(err)
	s.Require().Len(actual, 1)
	s.Require().Equal(compositePK, actual[0])

	err = os.RemoveAll(dir)
	s.Require().NoError(err)
}
