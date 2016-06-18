# reform [![GoDoc](https://godoc.org/github.com/go-reform/reform?status.svg)](https://godoc.org/github.com/go-reform/reform) [![Build Status](https://travis-ci.org/go-reform/reform.svg?branch=master)](https://travis-ci.org/go-reform/reform) [![Coverage Status](https://coveralls.io/repos/github/go-reform/reform/badge.svg?branch=master)](https://coveralls.io/github/go-reform/reform?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/go-reform/reform)](https://goreportcard.com/report/github.com/go-reform/reform)

A better ORM for Go.

It uses code generation (`go generate`), non-empty interfaces and initialization-time reflection
as opposed to type system sidestepping, `interface{}` and runtime reflection. It will be kept simple.


## Quickstart

1. Install it: `go get github.com/go-reform/reform/reform` (see about versioning below)
2. Define your first model in file `person.go`:

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
   in the same package with type `PersonTable` and methods on `Person`.
4. See [documentation](https://godoc.org/github.com/go-reform/reform) how to use it. Simple example:

    ```go
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
    ```


## Background

reform was born during summer 2014 out of frustrations with existing Go ORMs. All of them have a method
`Save(record interface{})` which can be used like this:

```go
orm.Save(User{Name: "gopher"})
orm.Save(&User{Name: "gopher"})
orm.Save(nil)
orm.Save("Batman!!")
```

Now you can say that last invocation is obviously invalid, and that it's not hard to make an ORM to accept both
first and second versions. But there are two problems:

1. Compiler can't check it. Method's signature in `godoc` will not tell us how to use it.
   We are essentially working against those tools by sidestepping type system.
2. First version is still invalid, since one would expect `Save()` method to set record's primary key after `INSERT`,
   but this change will be lost due to passing by value.

First proprietary version of reform was used in production even before `go generate` announcement.
This free and open-source version is the fourth milestone on the road to better and idiomatic API.


## Versioning policy

We are following Semantic Versioning, using [gopkg.in](http://gopkg.in) and filling a [changelog](CHANGELOG.md).

We use branch `v1-dev` (default on Github) for v1 development and tags `v1.Y.Z` for releases.
Code in that branch uses canonical import path `gopkg.in/reform.v1`.
Breaking changes will not be applied to v1.

We use `master` branch for what will became v2. Code there doesn't enforce canonical import path.


## Additional packages

* [github.com/AlekSi/pointer](https://github.com/AlekSi/pointer) is very useful for working with reform structs with pointers.
* [github.com/mc2soft/pq-types](https://github.com/mc2soft/pq-types) is a collection of PostgreSQL types, we use it with reform.


## Caveats

* There should be zero `pk` fields for Struct and exactly one `pk` field for Record.
* `pk` field can't be a pointer (`== nil` [doesn't work](https://golang.org/doc/faq#nil_error)).
* Database row can't have a Go's zero value (0, empty string, etc.) in primary key column.
