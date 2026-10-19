help:                    ## Display this help message.
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

# SHELL = go run .github/shell.go

REFORM_ENV_WAIT_TIMEOUT ?= 180

init:                    ## Install development tools.
	go mod verify
	cd tools && go mod verify
	cd tools && go generate -tags=tools -x

env-up:                  ## Start development environment.
	$(MAKE) env-up-detach
	trap 'docker compose stop' EXIT INT TERM; docker compose logs --follow

env-up-detach:           ## Start development environment in the background.
	docker compose up --wait --wait-timeout $(REFORM_ENV_WAIT_TIMEOUT) --force-recreate --renew-anon-volumes --remove-orphans || { \
		docker compose ps --all; \
		docker compose logs --no-color --tail=100; \
		exit 1; \
	}

env-down:                ## Stop development environment.
	docker compose down --volumes --remove-orphans

env-mysql:               ## Run mysql client.
	docker exec -ti reform_mysql mysql

test:                    ## Run all tests and gather coverage.
	make test-unit

	make postgres
	make pgx
	make mysql
	make mysql-traditional
	make sqlite3
	make mssql
	make sqlserver

	make merge-cover

test-unit:
	rm -f internal/test/models/*_reform.go
	rm -f reform-db/*_reform.go

	go install -v gopkg.in/reform.v1/reform
	go test -count=1 -race gopkg.in/reform.v1/parse
	go test -count=1 -covermode=count -coverprofile=parse.cover gopkg.in/reform.v1/parse
	go generate -v -x gopkg.in/reform.v1/internal/test/models
	go install -v gopkg.in/reform.v1/internal/test/models

	go generate -v -x gopkg.in/reform.v1/reform-db
	go install -v gopkg.in/reform.v1/reform-db

	go test -short -run '^Test' -count=1 -race ./...
	go vet ./...

test-db-init:
	# recreate and initialize database
	rm -f $(CURDIR)/reform-database.sqlite3
	-reform-db -db-driver="$(REFORM_TEST_DRIVER)" -db-source="$(REFORM_TEST_ADMIN_SOURCE)" -db-wait=15s exec \
		test/sql/$(REFORM_TEST_DATABASE)_drop.sql
	reform-db -db-driver="$(REFORM_TEST_DRIVER)" -db-source="$(REFORM_TEST_ADMIN_SOURCE)" exec \
		test/sql/$(REFORM_TEST_DATABASE)_create.sql
	reform-db -db-driver="$(REFORM_TEST_DRIVER)" -db-source="$(REFORM_TEST_INIT_SOURCE)" exec \
		test/sql/$(REFORM_TEST_DATABASE)_combined.tmp.sql

# run integration tests
test-db:
	# TODO remove that hack in reform 1.6
	# https://github.com/go-reform/reform/issues/151
	# https://github.com/go-reform/reform/issues/157
	cat \
		test/sql/$(REFORM_TEST_DATABASE)_init.sql \
		test/sql/data.sql \
		test/sql/$(REFORM_TEST_DATABASE)_data.sql \
		test/sql/$(REFORM_TEST_DATABASE)_set.sql \
		> test/sql/$(REFORM_TEST_DATABASE)_combined.tmp.sql

	make test-db-init

	# run reform-db tests
	go test -count=1 -race gopkg.in/reform.v1/reform-db
	go test -count=1 -covermode=count -coverprofile=$(REFORM_TEST_COVER)_reform-db.cover gopkg.in/reform.v1/reform-db

	# run main tests with -race
	go test -count=1 -race

	make test-db-init

	# run main tests with -cover
	go test -count=1 -covermode=count -coverprofile=$(REFORM_TEST_COVER)_reform.cover

# run integration tests for PostgreSQL (postgres driver)
postgres: export REFORM_TEST_DATABASE = postgres
postgres: export REFORM_TEST_DRIVER = postgres
postgres: export REFORM_TEST_ADMIN_SOURCE = postgres://postgres@127.0.0.1/template1?sslmode=disable
postgres: export REFORM_TEST_INIT_SOURCE = postgres://postgres@127.0.0.1/reform-database?sslmode=disable&TimeZone=UTC
postgres: export REFORM_TEST_SOURCE = postgres://postgres@127.0.0.1/reform-database?sslmode=disable&TimeZone=America/New_York
postgres: export REFORM_TEST_COVER=postgres
postgres:
	make test-db

# run integration tests for PostgreSQL (pgx driver)
pgx: export REFORM_TEST_DATABASE = postgres
pgx: export REFORM_TEST_DRIVER = pgx
pgx: export REFORM_TEST_ADMIN_SOURCE = postgres://postgres@127.0.0.1/template1?sslmode=disable
pgx: export REFORM_TEST_INIT_SOURCE = postgres://postgres@127.0.0.1/reform-database?sslmode=disable&TimeZone=UTC
pgx: export REFORM_TEST_SOURCE = postgres://postgres@127.0.0.1/reform-database?sslmode=disable&TimeZone=America/New_York
pgx: export REFORM_TEST_COVER=pgx
pgx:
	make test-db

# run integration tests for MySQL (ANSI SQL mode)
mysql: export REFORM_TEST_DATABASE = mysql
mysql: export REFORM_TEST_DRIVER = mysql
mysql: export REFORM_TEST_ADMIN_SOURCE = root@/mysql
mysql: export REFORM_TEST_INIT_SOURCE = root@/reform-database?parseTime=true&clientFoundRows=true&time_zone='UTC'&sql_mode='ANSI'&multiStatements=true
mysql: export REFORM_TEST_SOURCE = root@/reform-database?parseTime=true&clientFoundRows=true&time_zone='America%2FNew_York'&sql_mode='ANSI'
mysql: export REFORM_TEST_COVER=mysql
mysql:
	make test-db

# run integration tests for MySQL (traditional SQL mode + interpolateParams)
mysql-traditional: export REFORM_TEST_DATABASE = mysql
mysql-traditional: export REFORM_TEST_DRIVER = mysql
mysql-traditional: export REFORM_TEST_ADMIN_SOURCE = root@/mysql
mysql-traditional: export REFORM_TEST_INIT_SOURCE = root@/reform-database?parseTime=true&clientFoundRows=true&time_zone='UTC'&sql_mode='ANSI'&multiStatements=true
mysql-traditional: export REFORM_TEST_SOURCE = root@/reform-database?parseTime=true&clientFoundRows=true&time_zone='America%2FNew_York'&sql_mode='TRADITIONAL'&interpolateParams=true
mysql-traditional: export REFORM_TEST_COVER=mysql-traditional
mysql-traditional:
	make test-db

# run integration tests for SQLite3
sqlite3: export REFORM_TEST_DATABASE = sqlite3
sqlite3: export REFORM_TEST_DRIVER = sqlite3
sqlite3: export REFORM_TEST_ADMIN_SOURCE = $(CURDIR)/reform-database.sqlite3
sqlite3: export REFORM_TEST_INIT_SOURCE = $(CURDIR)/reform-database.sqlite3
sqlite3: export REFORM_TEST_SOURCE = $(CURDIR)/reform-database.sqlite3
sqlite3: export REFORM_TEST_COVER=sqlite3
sqlite3:
	make test-db

# run integration tests for SQL Server (mssql driver)
mssql: export REFORM_TEST_DATABASE = mssql
mssql: export REFORM_TEST_DRIVER = mssql
mssql: export REFORM_TEST_ADMIN_SOURCE = server=localhost;user id=sa;password=reform-password123
mssql: export REFORM_TEST_INIT_SOURCE = server=localhost;user id=sa;password=reform-password123;database=reform-database
mssql: export REFORM_TEST_SOURCE = server=localhost;user id=sa;password=reform-password123;database=reform-database
mssql: export REFORM_TEST_COVER=mssql
mssql:
	make test-db

# run integration tests for SQL Server (sqlserver driver)
sqlserver: export REFORM_TEST_DATABASE = mssql
sqlserver: export REFORM_TEST_DRIVER = sqlserver
sqlserver: export REFORM_TEST_ADMIN_SOURCE = server=localhost;user id=sa;password=reform-password123
sqlserver: export REFORM_TEST_INIT_SOURCE = server=localhost;user id=sa;password=reform-password123;database=reform-database
sqlserver: export REFORM_TEST_SOURCE = server=localhost;user id=sa;password=reform-password123;database=reform-database
sqlserver: export REFORM_TEST_COVER=sqlserver
sqlserver:
	make test-db

# Windows: run unit tests and integration tests for SQL Server (mssql driver)
win-mssql: REFORM_SQL_HOST ?= 127.0.0.1
win-mssql: REFORM_SQL_INSTANCE ?= SQLEXPRESS
win-mssql: export REFORM_TEST_DATABASE = mssql
win-mssql: export REFORM_TEST_DRIVER = mssql
win-mssql: export REFORM_TEST_ADMIN_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE)
win-mssql: export REFORM_TEST_INIT_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE);database=reform-database
win-mssql: export REFORM_TEST_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE);database=reform-database
win-mssql: export REFORM_TEST_COVER=win-mssql
win-mssql:
	make test-db

# Windows: run unit tests and integration tests for SQL Server (sqlserver driver)
win-sqlserver: REFORM_SQL_HOST ?= 127.0.0.1
win-sqlserver: REFORM_SQL_INSTANCE ?= SQLEXPRESS
win-sqlserver: export REFORM_TEST_DATABASE = mssql
win-sqlserver: export REFORM_TEST_DRIVER = sqlserver
win-sqlserver: export REFORM_TEST_ADMIN_SOURCE = sqlserver://$(REFORM_SQL_HOST)/$(REFORM_SQL_INSTANCE)
win-sqlserver: export REFORM_TEST_INIT_SOURCE = sqlserver://$(REFORM_SQL_HOST)/$(REFORM_SQL_INSTANCE)?database=reform-database
win-sqlserver: export REFORM_TEST_SOURCE = sqlserver://$(REFORM_SQL_HOST)/$(REFORM_SQL_INSTANCE)?database=reform-database
win-sqlserver: export REFORM_TEST_COVER=win-sqlserver
win-sqlserver:
	make test-db

merge-cover:
	bin/gocoverutil -coverprofile=coverage.txt merge *.cover
	rm -f *.cover

lint:                    ## Run linters.
	# run required linters for all code
	bin/golangci-lint run --config=.golangci-required.yml
	bin/golangci-lint fmt --config=.golangci-required.yml --diff

	# run optional linters for new code only
	bin/golangci-lint run --new --issues-exit-code=0
	bin/go-consistent -pedantic ./... | bin/reviewdog -f=go-consistent -diff='git diff HEAD^'

ci-check-changes:
	go mod tidy -diff
	cd tools && go mod tidy -diff
	git diff --exit-code

.PHONY: help init env-up env-up-detach env-down env-mysql test test-unit test-db-init test-db \
	postgres pgx mysql mysql-traditional sqlite3 mssql sqlserver win-mssql win-sqlserver merge-cover lint ci-check-changes
