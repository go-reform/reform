package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/reform.v1/parse"
)

func (s *ReformDBSuite) TestInit() {
	good, err := parse.File("../internal/test/models/good.go")
	s.Require().NoError(err)
	s.Require().Len(good, 5)

	dir, err := ioutil.TempDir("", "TestInit")
	s.Require().NoError(err)
	s.T().Log(dir)

	cmdInit(s.db, dir)

	ff := filepath.Join(dir, "projects.go")
	actual, err := parse.File(ff)
	s.Require().NoError(err)

	s.Require().Equal(good[1], actual[0])

	err = os.RemoveAll(dir)
	s.Require().NoError(err)
}
