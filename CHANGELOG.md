# Changelog

## v1.2.1 (2016-09-14, https://github.com/go-reform/reform/milestones/v1.2.1)

* `reform` tool now correctly handles non-exported types.
* Querier.Insert now correctly INSERTs records with set non-integer primary keys, even if
  dialect uses LastInsertId (MySQL, SQLite3).

## v1.2.0 (2016-08-10, https://github.com/go-reform/reform/milestones/v1.2.0)

* Added support for Microsoft SQL Server. Huge thanks to [Aleksey Martynov](https://github.com/AlekseyMartynov).
* Added Querier.InsertColumns.
* Querier.Insert now correctly handles records with only primary key column.

## v1.1.2 (2016-07-20, https://github.com/go-reform/reform/milestones/v1.1.2)

* `reform` tool now correctly ignores type information when it's not used.
  This allows one to have fields of any custom types. The only exception is primary key fields,
  which are restricted to basic types (numbers and strings).
* Package `gopkg.in/reform.v1/parse` is explicitly documented as internal.
  (It's wasn't really possible to use it.)

## v1.1.1 (2016-07-05, https://github.com/go-reform/reform/milestones/v1.1.1)

* Querier.UpdateColumns no longer allows to update primary key column. This behavior was allowed,
  but did not make any sense.
* `reform` tool now correctly handles pointers to custom types and slices.

## v1.1.0 (2016-07-01, https://github.com/go-reform/reform/milestones/v1.1.0)

* Added Querier.InsertMulti.
* Added DBInterface, TXInterface, NewDBFromInterface, NewTXFromInterface.

## v1.0.0 (2016-06-22)

* Moved to https://github.com/go-reform/reform repository.
* Changed canonical import path.
* Added versioning policy.
