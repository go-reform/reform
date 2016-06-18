package reform_test

import (
	"fmt"
	"log"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/go-reform/reform"
	. "github.com/go-reform/reform/internal/test/models"
)

func Example() {
	// save record (performs INSERT or UPDATE)
	person := &Person{
		Name:  "Alexey Palazhchenko",
		Email: pointer.ToString("alexey.palazhchenko@gmail.com"),
	}
	if err := DB.Save(person); err != nil {
		log.Fatal(err)
	}

	// ID is filled by Save
	person2, err := DB.FindByPrimaryKeyFrom(PersonTable, person.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(person2.(*Person).Name)

	// delete record
	if err = DB.Delete(person); err != nil {
		log.Fatal(err)
	}

	// find records by IDs
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
