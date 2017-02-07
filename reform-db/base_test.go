package main

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/internal"
)

type ReformDBSuite struct {
	suite.Suite
	db *reform.DB
}

func TestReformDBSuite(t *testing.T) {
	suite.Run(t, new(ReformDBSuite))
}

func (s *ReformDBSuite) SetupSuite() {
	logger = internal.NewLogger("reform-db-test: ", true)

	s.db = internal.ConnectToTestDB()
	s.db.Querier = s.db.WithTag("reform-db-test")
}

func (s *ReformDBSuite) SetupTest() {
	pl := reform.NewPrintfLogger(s.T().Logf)
	pl.LogTypes = true
	s.db.Logger = pl
}
