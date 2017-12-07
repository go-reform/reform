# reform
[![Release](https://github-release-version.herokuapp.com/github/go-reform/reform/release.svg?style=flat)](https://github.com/go-reform/reform/releases/latest)
[![GoDoc](https://godoc.org/gopkg.in/reform.v1?status.svg)](https://godoc.org/gopkg.in/reform.v1)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-reform/reform?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Travis CI Build Status](https://travis-ci.org/go-reform/reform.svg?branch=v1-stable)](https://travis-ci.org/go-reform/reform)
[![AppVeyor Build status](https://ci.appveyor.com/api/projects/status/kbkyjmic461xa7b3/branch/v1-stable?svg=true)](https://ci.appveyor.com/project/AlekSi/reform/branch/v1-stable)
[![Coverage Report](https://codecov.io/gh/go-reform/reform/branch/v1-stable/graph/badge.svg)](https://codecov.io/gh/go-reform/reform)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/reform.v1)](https://goreportcard.com/report/gopkg.in/reform.v1)

<a href="https://en.wikipedia.org/wiki/Peter_the_Great"><img align="right" alt="Reform gopher logo" title="Peter the Reformer" src=".github/reform.png"></a>

A better ORM for Go and `database/sql`.

It uses non-empty interfaces, code generation (`go generate`), and initialization-time reflection
as opposed to `interface{}`, type system sidestepping, and runtime reflection. It will be kept simple.

Supported SQL dialects:

| RDBMS                | Library and drivers                                                                                 | Tested with
| -----                | -------------------                                                                                 | -----------
| PostgreSQL           | [github.com/lib/pq](https://github.com/lib/pq) (`postgres`)                                         | All [supported](https://www.postgresql.org/support/versioning/) versions.
| MySQL                | [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) (`mysql`)                  | All [supported](https://www.mysql.com/support/supportedplatforms/database.html) versions.
| SQLite3              | [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) (`sqlite3`)                      |
| Microsoft SQL Server | [github.com/denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb) (`mssql`, `sqlserver`) | Windows: SQL2008R2SP2, SQL2012SP1, SQL2014, SQL2016. Linux: [`microsoft/mssql-server-linux:latest` Docker image](https://hub.docker.com/r/microsoft/mssql-server-linux/).

Note that for MySQL [`clientFoundRows=true`](https://github.com/go-sql-driver/mysql#clientfoundrows) flag is required.

## Quickstart

1. Make sure you are using Go 1.7+. Install or update `reform` package, `reform` and `reform-db` commands
   (see about versioning below):

    ```
    go get -u gopkg.in/reform.v1/...
    ```

2. Use `reform-db` command to generate models for your existing database schema. For example:
    ```
    reform-db -db-driver=sqlite3 -db-source=example.sqlite3 init
    ```

3. Update generated models or write your own – `struct` representing a table or view row. For example,
   store this in file `person.go`:

    ```go
    //go:generate reform

    //reform:people
	type Person struct {
		ID        int32      `reform:"id,pk"`
		Name      string     `reform:"name"`
		Email     *string    `reform:"email"`
		CreatedAt time.Time  `reform:"created_at"`
		UpdatedAt *time.Time `reform:"updated_at"`
	}
    ```

    Magic comment `//reform:people` links this model to `people` table or view in SQL database.
    The first value in field's `reform` tag is a column name. `pk` marks primary key.
    Use value `-` or omit tag completely to skip a field.
    Use pointers for nullable fields.

4. Run `reform [package or directory]` or `go generate [package or file]`. This will create `person_reform.go`
   in the same package with type `PersonTable` and methods on `Person`.
5. See [documentation](https://godoc.org/gopkg.in/reform.v1) how to use it. Simple example:

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

We are following [Semantic Versioning](http://semver.org/spec/v2.0.0.html),
using [gopkg.in](https://gopkg.in) and filling a [changelog](CHANGELOG.md).

We use branch `v1-stable` (default on Github) for v1 development and tags `v1.Y.Z` for releases.
All v1 releases are SemVer-compatible, breaking changes will not be applied.
Canonical import path is `gopkg.in/reform.v1`.
`go get -u gopkg.in/reform.v1/reform` will install latest released version.
To install not yet released v1 version one can do checkout manually while preserving import path:
```
cd $GOPATH/src/gopkg.in/reform.v1
git fetch
git checkout origin/v1-stable
go install -v gopkg.in/reform.v1/reform
```

Branch `v2-unstable` is used for v2 development. It doesn't have any releases yet, and no compatibility is guaranteed.
Canonical import path is `gopkg.in/reform.v2-unstable`.


## Additional packages

* [github.com/AlekSi/pointer](https://github.com/AlekSi/pointer) is very useful for working with reform structs with pointers.


## Caveats and limitations

* There should be zero `pk` fields for Struct and exactly one `pk` field for Record.
  Composite primary keys are not supported ([#114](https://github.com/go-reform/reform/issues/114)).
* `pk` field can't be a pointer (`== nil` [doesn't work](https://golang.org/doc/faq#nil_error)).
* Database row can't have a Go's zero value (0, empty string, etc.) in primary key column.


## License

Code is covered by standard MIT-style license. Copyright (c) 2016-2017 Alexey Palazhchenko.
See [LICENSE](LICENSE) for details. Note that generated code is covered by the terms of your choice.

The reform gopher was drawn by Natalya Glebova. Please use it only as reform logo.
It is based on the original design by Renée French, released under [Creative Commons Attribution 3.0 USA license](https://creativecommons.org/licenses/by/3.0/).


## Contributing

See [Contributing Guidelines](.github/CONTRIBUTING.md).
