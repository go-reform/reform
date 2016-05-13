# reform [![GoDoc](https://godoc.org/github.com/AlekSi/reform?status.svg)](https://godoc.org/github.com/AlekSi/reform)

A better ORM for Go.

It uses non-empty interfaces, code generation (`go generate`) and initialization-time reflection
as opposed to `interface{}` and runtime reflection. It

## Quickstart

1. Install it: `go get github.com/AlekSi/reform/reform` (see about versioning below)
2. Define your first model in file person.go:

    ```go
    //go:generate reform

    //reform:people
	Person struct {
		ID        int32      `reform:"id,pk"`
		Name      string     `reform:"name"`
		Email     *string    `reform:"email"`
		CreatedAt time.Time  `reform:"created_at"`
		UpdatedAt *time.Time `reform:"updated_at"`
	}
    ```

    Magic comment `//reform:people` links this model to `people` table or view in SQL database.
    First value in `reform` tag is a column name. `pk` marks primary key.
    Use pointers for nullable fields.

3. Run `reform [package or directory]` or `go generate [package or file]`. This will create `person_reform.go`
   in the same package.
4. See [documentation](https://godoc.org/github.com/AlekSi/reform) how to use it. Simple example:

    ```go
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
    ```

## Versioning

TODO

## Additional packages

* [github.com/AlekSi/pointer](https://github.com/AlekSi/pointer) is very useful for working with reform structs with pointers.
* [github.com/mc2soft/pq-types](https://github.com/mc2soft/pq-types) is a collection of PostgreSQL types, we use it with reform.

## Caveats

* There should be zero `pk` fields for Struct and exactly one `pk` field for Record.
* `pk` field can't be a pointer (`== nil` [doesn't work](https://golang.org/doc/faq#nil_error)).
Database row can't have a Go's zero value in primary key column.
