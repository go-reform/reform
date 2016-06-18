package reform_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/go-reform/reform"
	"github.com/go-reform/reform/dialects/postgresql"
	. "github.com/go-reform/reform/internal/test/models"
)

func Example() {
	// Use reform.NewDB to create DB.

	// Save record (performs INSERT or UPDATE).
	person := &Person{
		Name:  "Alexey Palazhchenko",
		Email: pointer.ToString("alexey.palazhchenko@gmail.com"),
	}
	if err := DB.Save(person); err != nil {
		log.Fatal(err)
	}

	// ID is filled by Save.
	person2, err := DB.FindByPrimaryKeyFrom(PersonTable, person.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(person2.(*Person).Name)

	// Delete record.
	if err = DB.Delete(person); err != nil {
		log.Fatal(err)
	}

	// Find records by IDs.
	persons, err := DB.FindAllFrom(PersonTable, "id", 1, 2)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range persons {
		fmt.Println(p)
	}
	// Output:
	// Alexey Palazhchenko
	// ID: 1 (int32), Name: `Denis Mills` (string), Email: <nil> (*string), CreatedAt: 2009-11-10 23:00:00 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
	// ID: 2 (int32), Name: `Garrick Muller` (string), Email: `muller_garrick@example.com` (*string), CreatedAt: 2009-12-12 12:34:56 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
}

func ExampleNewDB() {
	// Get *sql.DB as usual. PostgreSQL example:
	conn, err := sql.Open("postgres", "postgres://localhost:5432/reform-test?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Use new *log.Logger for logging.
	logger := log.New(os.Stderr, "SQL: ", log.Flags())

	// Create *reform.DB instance with simple logger.
	// Any Printf-like function (fmt.Printf, log.Printf, testing.T.Logf, etc) can be used with NewPrintfLogger.
	// Change dialect for other databases.
	_ = reform.NewDB(conn, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))
}

func ExampleQuerier_SelectRows() {
	tail := fmt.Sprintf("WHERE created_at < %s ORDER BY id", DB.Placeholder(1))
	y2010 := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	rows, err := DB.SelectRows(PersonTable, tail, y2010)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for {
		var person Person
		err = DB.NextRow(&person, rows)
		if err != nil {
			break
		}
		fmt.Println(person)
	}
	if err != reform.ErrNoRows {
		log.Fatal(err)
	}
	// Output:
	// ID: 1 (int32), Name: `Denis Mills` (string), Email: <nil> (*string), CreatedAt: 2009-11-10 23:00:00 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
	// ID: 2 (int32), Name: `Garrick Muller` (string), Email: `muller_garrick@example.com` (*string), CreatedAt: 2009-12-12 12:34:56 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
}

func ExampleQuerier_SelectOneTo() {
	var person Person
	tail := fmt.Sprintf("WHERE created_at < %s ORDER BY id", DB.Placeholder(1))
	y2010 := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	err := DB.SelectOneTo(&person, tail, y2010)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(person)
	// Output:
	// ID: 1 (int32), Name: `Denis Mills` (string), Email: <nil> (*string), CreatedAt: 2009-11-10 23:00:00 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
}
