package reform_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/enodata/faker"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	. "gopkg.in/reform.v1/internal/test/models"
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
	// ID: 1 (int32), GroupID: 65534 (*int32), Name: `Denis Mills` (string), Email: <nil> (*string), CreatedAt: 2009-11-10 23:00:00 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
	// ID: 2 (int32), GroupID: 65534 (*int32), Name: `Garrick Muller` (string), Email: `muller_garrick@example.com` (*string), CreatedAt: 2009-12-12 12:34:56 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
}

func ExampleNewDB() {
	// Get *sql.DB as usual. PostgreSQL example:
	conn, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=disable")
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

func ExampleQuerier_WithTag() {
	id := "baron"
	person, err := DB.WithTag("GetProject:%v", id).FindByPrimaryKeyFrom(ProjectTable, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(person)
	// Output:
	// Name: `Vicious Baron` (string), ID: `baron` (string), Start: 2014-06-01 00:00:00 +0000 UTC (time.Time), End: 2016-02-21 00:00:00 +0000 UTC (*time.Time)
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
		if err = DB.NextRow(&person, rows); err != nil {
			break
		}
		fmt.Println(person)
	}
	if err != reform.ErrNoRows {
		log.Fatal(err)
	}
	// Output:
	// ID: 1 (int32), GroupID: 65534 (*int32), Name: `Denis Mills` (string), Email: <nil> (*string), CreatedAt: 2009-11-10 23:00:00 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
	// ID: 2 (int32), GroupID: 65534 (*int32), Name: `Garrick Muller` (string), Email: `muller_garrick@example.com` (*string), CreatedAt: 2009-12-12 12:34:56 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
}

func ExampleQuerier_SelectOneTo() {
	var person Person
	tail := fmt.Sprintf("WHERE created_at < %s ORDER BY id", DB.Placeholder(1))
	y2010 := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := DB.SelectOneTo(&person, tail, y2010); err != nil {
		log.Fatal(err)
	}
	fmt.Println(person)
	// Output:
	// ID: 1 (int32), GroupID: 65534 (*int32), Name: `Denis Mills` (string), Email: <nil> (*string), CreatedAt: 2009-11-10 23:00:00 +0000 UTC (time.Time), UpdatedAt: <nil> (*time.Time)
}

var persons = []reform.Struct{
	&Person{
		Name:  "Alexey Palazhchenko",
		Email: pointer.ToString("alexey.palazhchenko@gmail.com"),
	},
	&Person{
		Name:  faker.Name().Name(),
		Email: pointer.ToString(faker.Internet().Email()),
	},
	&Person{
		Name:  faker.Name().Name(),
		Email: pointer.ToString(faker.Internet().Email()),
	},
	&Person{
		Name:  faker.Name().Name(),
		Email: pointer.ToString(faker.Internet().Email()),
	},
	&Person{
		Name:  faker.Name().Name(),
		Email: pointer.ToString(faker.Internet().Email()),
	},
}

func ExampleQuerier_InsertMulti() {
	// insert up to 3 structs at once
	const batchSize = 3
	for i := 0; i < len(persons)/batchSize+1; i++ {
		low := i * batchSize
		high := (i + 1) * batchSize
		if high > len(persons) {
			high = len(persons)
		}
		batch := persons[low:high]

		if err := DB.InsertMulti(batch...); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted %d persons\n", len(batch))
	}

	// note that ID is not filled
	fmt.Println(persons[0].(*Person).ID, persons[0].(*Person).Name)
	// Output:
	// Inserted 3 persons
	// Inserted 2 persons
	// 0 Alexey Palazhchenko
}

func ExampleQuerier_Query() {
	columns := DB.QualifiedColumns(PersonTable)
	columns = append(columns, DB.QualifiedColumns(PersonProjectView)...)
	columns = append(columns, DB.QualifiedColumns(ProjectTable)...)
	query := fmt.Sprintf(`
		SELECT %s
			FROM people
			INNER JOIN person_project ON people.id = person_project.person_id
			INNER JOIN projects ON person_project.project_id = projects.id
			ORDER BY person_id, project_id;
	`, strings.Join(columns, ", "))
	rows, err := DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var person Person
		var personProject PersonProject
		var project Project
		pointers := person.Pointers()
		pointers = append(pointers, personProject.Pointers()...)
		pointers = append(pointers, project.Pointers()...)
		if err = rows.Scan(pointers...); err != nil {
			log.Print(err)
		}
		if err = person.AfterFind(); err != nil {
			log.Fatal(err)
		}
		if err = project.AfterFind(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s - %s\n", person.Name, project.Name)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	// Output:
	// Noble Schumm - Vicious Baron
	// Elfrieda Abbott - Vicious Baron
	// Elfrieda Abbott - Thirsty Queen
	// Elfrieda Abbott - Vicious Baron
	// Elfrieda Abbott - Thirsty Queen
	// Elfrieda Abbott - Kosher Traveler
}
