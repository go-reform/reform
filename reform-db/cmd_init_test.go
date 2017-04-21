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
	s.Require().Len(good, 5)

	people := good[0]
	projects := good[1]
	personProject := good[2]
	idOnly := good[4]

	// patch difference we don't handle
	people.Type = strings.Replace(people.Type, "Person", "People", -1)
	projects.Type = strings.Replace(projects.Type, "Project", "Projects", -1)
	if s.db.Dialect == sqlite3.Dialect {
		people.Fields[0].Type = strings.Replace(people.Fields[0].Type, "int32", "int64", -1)
		people.Fields[1].Type = strings.Replace(people.Fields[1].Type, "int32", "int64", -1)
		personProject.Fields[0].Type = strings.Replace(personProject.Fields[0].Type, "int32", "int64", -1)
		idOnly.Fields[0].Type = strings.Replace(idOnly.Fields[0].Type, "int32", "int64", -1)
	}

	dir, err := ioutil.TempDir("", "ReformDBTestInit")
	s.Require().NoError(err)
	s.T().Log(dir)

	cmdInit(s.db, dir)

	fis, err := ioutil.ReadDir(dir)
	s.Require().NoError(err)
	s.Require().Len(fis, 4)

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

	err = os.RemoveAll(dir)
	s.Require().NoError(err)
}
