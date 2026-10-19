# Contributing Guidelines

## Golden rule

Speak up _before_ writing code.
Comment on existing issue or create a new one.
Discuss what you want to implement _before_ implementing it.

## Getting code

Use Go 1.27.2 for development (Go 1.26 is also supported), a C compiler for SQLite and race tests,
and Docker with a current Compose plugin. Add `$(go env GOPATH)/bin` to your `PATH` so the generated
`reform` and `reform-db` commands can be found.

Fork repository on GitHub and clone the source code:

```
git clone git@github.com:<your name>/reform.git
cd reform
make init
```

Please read the "Versioning and branching policy" section in README,
and send pull requests to the right branch.

## Makefile targets

* Run `make` without arguments to see all Makefile targets.
* Run `make env-up` or `make env-up-detach` to start databases with Docker Compose.
* Run `make init` to install development tools.
* Run `make test` to run all tests.
* Run `make lint` to run linters.

Database startup waits up to 180 seconds and prints service status and logs on failure.
Use `make env-up-detach REFORM_ENV_WAIT_TIMEOUT=300` to change the readiness timeout.
SQL Server uses a `linux/amd64` image; on Apple Silicon, Docker must support amd64 emulation.
Override `REFORM_POSTGRES_IMAGE`, `REFORM_MYSQL_IMAGE`, or `REFORM_MSSQL_IMAGE` to test other versions.

Run database-free tests with `go test -short -run '^Test' -count=1 -race ./...`.
The complete `make test` run also exercises database-dependent examples and requires the running databases.
