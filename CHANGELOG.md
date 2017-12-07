# Changelog

## v1.3.1 (2017-12-07, https://github.com/go-reform/reform/milestones/v1.3.1)

* No user-visible changes.
* Major changes in CI and development environment.

## v1.3.0 (2017-12-01, https://github.com/go-reform/reform/milestones/v1.3.0)

* Go 1.7+ is now required.
* Added `reform-db` command.
  * `init` subcommand may be used to generate Go model files for existing database schema.
  * `query` and `exec` subcommands may be used for accessing a database.
* Fields with `reform` tag with value `"-"` are ignored now (just like with value `""` and without tag at all).
* Added [`ErrTxDone`](https://godoc.org/gopkg.in/reform.v1#pkg-variables).
* Added [`DB.DBInterface`](https://godoc.org/gopkg.in/reform.v1#DB.DBInterface).
* Added [`Querier.UpdateView`](https://godoc.org/gopkg.in/reform.v1#Querier.UpdateView).
* `reform` command with `-gofmt=false` flag still formats generated sources with go/format package, without invoking `gofmt`.
  Thanks to [João Pereira](https://github.com/joaodrp).
* Added support for `sqlserver` variant of [github.com/denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb) driver.
* Added support for Microsoft SQL Server for Linux.
* We now have a logo! Huge thanks to Natalya Glebova for making it.

## v1.2.1 (2016-09-14, https://github.com/go-reform/reform/milestones/v1.2.1)

* `reform` command now correctly handles non-exported types.
* [`Querier.Insert`](https://godoc.org/gopkg.in/reform.v1#Querier.Insert) now correctly INSERTs records with set
  non-integer primary keys, even if dialect uses LastInsertId (MySQL, SQLite3).

## v1.2.0 (2016-08-10, https://github.com/go-reform/reform/milestones/v1.2.0)

* Added support for Microsoft SQL Server. Huge thanks to [Aleksey Martynov](https://github.com/AlekseyMartynov).
* Added [`Querier.InsertColumns`](https://godoc.org/gopkg.in/reform.v1#Querier.InsertColumns).
* [`Querier.Insert`](https://godoc.org/gopkg.in/reform.v1#Querier.Insert) now correctly handles records with only primary key column.

## v1.1.2 (2016-07-20, https://github.com/go-reform/reform/milestones/v1.1.2)

* `reform` command now correctly ignores type information when it's not used.
  This allows one to have fields of any custom types. The only exception is primary key fields,
  which are restricted to basic types (numbers and strings).
* Package [`gopkg.in/reform.v1/parse`](https://godoc.org/gopkg.in/reform.v1/parse) is explicitly documented as internal.
  (It's wasn't really possible to use it.)

## v1.1.1 (2016-07-05, https://github.com/go-reform/reform/milestones/v1.1.1)

* [`Querier.UpdateColumns`](https://godoc.org/gopkg.in/reform.v1#Querier.UpdateColumns) no longer allows to update
  primary key column. This behavior was allowed, but did not make any sense.
* `reform` command now correctly handles pointers to custom types and slices.

## v1.1.0 (2016-07-01, https://github.com/go-reform/reform/milestones/v1.1.0)

* Added [`Querier.InsertMulti`](https://godoc.org/gopkg.in/reform.v1#Querier.InsertMulti).
* Added [`DBInterface`](https://godoc.org/gopkg.in/reform.v1#DBInterface),
  [`TXInterface`](https://godoc.org/gopkg.in/reform.v1#TXInterface),
  [`NewDBFromInterface`](https://godoc.org/gopkg.in/reform.v1#NewDBFromInterface),
  [`NewTXFromInterface`](https://godoc.org/gopkg.in/reform.v1#NewTXFromInterface).

## v1.0.0 (2016-06-22)

* Moved to https://github.com/go-reform/reform repository.
* Changed canonical import path.
* Added versioning policy.
