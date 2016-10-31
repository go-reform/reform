package reform_test

import (
	"database/sql"
	"log"
	"os"

	"github.com/AlekSi/pointer"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	. "gopkg.in/reform.v1/internal/test/models"
)

func initMSDB() *reform.DB {
	driver := os.Getenv("REFORM_TEST_DRIVER")
	masterSource := os.Getenv("REFORM_TEST_SOURCE_MASTER")
	slaveSource := os.Getenv("REFORM_TEST_SOURCE_SLAVE")
	log.Printf("driver = %q, master_source = %q slave_source = %q", driver, masterSource, slaveSource)
	if driver == "" || masterSource == "" || slaveSource == "" {
		log.Fatal("no driver or source or slaveSource, set REFORM_TEST_DRIVER , REFORM_TEST_SOURCE_MASTER, REFORM_TEST_SOURCE_SLAVE")
	}

	db, err := sql.Open(driver, masterSource)
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
	slavedb, err := sql.Open(driver, slaveSource)
	if err != nil {
		log.Fatal(err)
	}
	slavedb.SetMaxIdleConns(1)
	slavedb.SetMaxOpenConns(1)
	slavedb.SetConnMaxLifetime(-1)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return reform.NewDB(db, slavedb, postgresql.Dialect, nil, true)
}

func (s *ReformSuite) TestReadNotInTransactionOnSlave() {
	MSDB := initMSDB()
	if _, err := MSDB.Querier.DeleteFrom(PersonTable, ""); err != nil {
		log.Fatal(err)
	}

	// When not in transaction read from slave
	person := &Person{
		Name:  "Alexey Palazhchenko",
		Email: pointer.ToString("alexey.palazhchenko@gmail.com"),
	}
	if err := MSDB.Save(person); err != nil {
		log.Fatal(err)
	}
	person.Name = "Roman Potekhin"
	if err := MSDB.Save(person); err != nil {
		log.Fatal(err)
	}
	if res, err := MSDB.Querier.SelectAllFrom(PersonTable, ""); err != nil {
		log.Fatal(err)
	} else {
		s.Equal(len(res), 0)
	}
}

func (s *ReformSuite) TestWriteReadInTransactionOnMaster() {
	MSDB := initMSDB()
	if _, err := MSDB.Querier.DeleteFrom(PersonTable, ""); err != nil {
		log.Fatal(err)
	}

	// When in transaction write and read from master
	person := &Person{
		Name:  "Alexey Palazhchenko",
		Email: pointer.ToString("alexey.palazhchenko@gmail.com"),
	}

	// Transaction func
	MSDB.InTransaction(func(tx *reform.TX) error {
		if err := tx.Save(person); err != nil {
			return err
		}
		if res, err := tx.SelectAllFrom(PersonTable, ""); err != nil {
			return err
		} else {
			s.Equal(len(res), 1)
		}
		return nil
	})

	if _, err := MSDB.Querier.DeleteFrom(PersonTable, ""); err != nil {
		log.Fatal(err)
	}

	// pure tx
	tx, err := MSDB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	if err := tx.Save(person); err != nil {
		log.Fatal(err)
	}
	if res, err := tx.SelectAllFrom(PersonTable, ""); err != nil {
		log.Fatal(err)
	} else {
		s.Equal(len(res), 1)
	}
	tx.Rollback()
}
