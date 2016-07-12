# Changelog

## v1.2.0 (not yet released)

* Added Querier.InsertColumns.

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
