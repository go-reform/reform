all: test-dc check

# extra flags like -v
REFORM_TEST_FLAGS ?=

# SHELL = go run .github/shell.go

# install dependencies
deps:
	go get -u github.com/lib/pq
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/mattn/go-sqlite3
	go get -u github.com/denisenkom/go-mssqldb

	go get -u github.com/AlekSi/pointer
	go get -u github.com/stretchr/testify/...
	go get -u github.com/enodata/faker
	go get -u gopkg.in/alecthomas/gometalinter.v1
	go get -u github.com/AlekSi/gocoverutil

	gometalinter.v1 --install

# run all linters
check:
	-gometalinter.v1 ./... --deadline=180s --severity=vet:error

# run unit tests, generate models, install tools
test:
	rm -f *.cover coverage.txt
	rm -f internal/test/models/*_reform.go
	rm -f reform-db/*_reform.go

	go install -v gopkg.in/reform.v1/reform
	go test $(REFORM_TEST_FLAGS) -covermode=count -coverprofile=parse.cover gopkg.in/reform.v1/parse
	go generate -v -x gopkg.in/reform.v1/internal/test/models
	go install -v gopkg.in/reform.v1/internal/test/models

	go generate -v -x gopkg.in/reform.v1/reform-db
	go install -v gopkg.in/reform.v1/reform-db

# initialize database and run integration tests
test-db:
	-reform-db -db-driver="$(REFORM_DRIVER)" -db-source="$(REFORM_ROOT_SOURCE)" -db-wait=15s exec \
		internal/test/sql/$(REFORM_DATABASE)_drop.sql
	reform-db -db-driver="$(REFORM_DRIVER)" -db-source="$(REFORM_ROOT_SOURCE)" exec \
		internal/test/sql/$(REFORM_DATABASE)_create.sql
	reform-db -db-driver="$(REFORM_DRIVER)" -db-source="$(REFORM_INIT_SOURCE)" exec \
		internal/test/sql/$(REFORM_DATABASE)_init.sql \
		internal/test/sql/data.sql \
		internal/test/sql/$(REFORM_DATABASE)_data.sql \
		internal/test/sql/$(REFORM_DATABASE)_set.sql
	go test $(REFORM_TEST_FLAGS) -covermode=count -coverprofile=reform-db.cover gopkg.in/reform.v1/reform-db
	go test $(REFORM_TEST_FLAGS) -covermode=count -coverprofile=reform.cover
	gocoverutil -coverprofile=coverage.txt merge *.cover
	rm -f *.cover

# run all integration tests with Docker Compose
test-dc:
	go run .github/test-dc.go test

# run unit tests and integration tests for PostgreSQL
postgres: export REFORM_DATABASE = postgres
postgres: export REFORM_DRIVER = postgres
postgres: export REFORM_ROOT_SOURCE = postgres://postgres@127.0.0.1/template1?sslmode=disable
postgres: export REFORM_INIT_SOURCE = postgres://postgres@127.0.0.1/reform-database?sslmode=disable&TimeZone=UTC
postgres: export REFORM_TEST_SOURCE = postgres://postgres@127.0.0.1/reform-database?sslmode=disable&TimeZone=America/New_York
postgres: test
	make test-db

# run unit tests and integration tests for MySQL (ANSI SQL mode)
mysql: export REFORM_DATABASE = mysql
mysql: export REFORM_DRIVER = mysql
mysql: export REFORM_ROOT_SOURCE = root@/mysql
mysql: export REFORM_INIT_SOURCE = root@/reform-database?parseTime=true&clientFoundRows=true&time_zone='UTC'&sql_mode='ANSI'&multiStatements=true
mysql: export REFORM_TEST_SOURCE = root@/reform-database?parseTime=true&clientFoundRows=true&time_zone='America%2FNew_York'&sql_mode='ANSI'
mysql: test
	make test-db

# run unit tests and integration tests for MySQL (traditional SQL mode + interpolateParams)
mysql-traditional: export REFORM_DATABASE = mysql
mysql-traditional: export REFORM_DRIVER = mysql
mysql-traditional: export REFORM_ROOT_SOURCE = root@/mysql
mysql-traditional: export REFORM_INIT_SOURCE = root@/reform-database?parseTime=true&clientFoundRows=true&time_zone='UTC'&sql_mode='ANSI'&multiStatements=true
mysql-traditional: export REFORM_TEST_SOURCE = root@/reform-database?parseTime=true&clientFoundRows=true&time_zone='America%2FNew_York'&sql_mode='TRADITIONAL'&interpolateParams=true
mysql-traditional: test
	make test-db

# run unit tests and integration tests for SQLite3
sqlite3: export REFORM_DATABASE = sqlite3
sqlite3: export REFORM_DRIVER = sqlite3
sqlite3: export REFORM_ROOT_SOURCE = /tmp/reform-database.sqlite3
sqlite3: export REFORM_INIT_SOURCE = /tmp/reform-database.sqlite3
sqlite3: export REFORM_TEST_SOURCE = /tmp/reform-database.sqlite3
sqlite3: test
	rm -f /tmp/reform-database.sqlite3
	make test-db

# run unit tests and integration tests for SQL Server (mssql driver)
mssql: export REFORM_DATABASE = mssql
mssql: export REFORM_DRIVER = mssql
mssql: export REFORM_ROOT_SOURCE = server=localhost;user id=sa;password=reform-password123
mssql: export REFORM_INIT_SOURCE = server=localhost;user id=sa;password=reform-password123;database=reform-database
mssql: export REFORM_TEST_SOURCE = server=localhost;user id=sa;password=reform-password123;database=reform-database
mssql: test
	make test-db

# run unit tests and integration tests for SQL Server (sqlserver driver)
sqlserver: export REFORM_DATABASE = mssql
sqlserver: export REFORM_DRIVER = sqlserver
sqlserver: export REFORM_ROOT_SOURCE = server=localhost;user id=sa;password=reform-password123
sqlserver: export REFORM_INIT_SOURCE = server=localhost;user id=sa;password=reform-password123;database=reform-database
sqlserver: export REFORM_TEST_SOURCE = server=localhost;user id=sa;password=reform-password123;database=reform-database
sqlserver: test
	make test-db

# Windows: run unit tests and integration tests for SQL Server (mssql driver)
win-mssql: REFORM_SQL_HOST ?= 127.0.0.1
win-mssql: REFORM_SQL_INSTANCE ?= SQLEXPRESS
win-mssql: export REFORM_DATABASE = mssql
win-mssql: export REFORM_DRIVER = mssql
win-mssql: export REFORM_ROOT_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE)
win-mssql: export REFORM_INIT_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE);database=reform-database
win-mssql: export REFORM_TEST_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE);database=reform-database
win-mssql: test
	mingw32-make test-db

# Windows: run unit tests and integration tests for SQL Server (sqlserver driver)
win-sqlserver: REFORM_SQL_HOST ?= 127.0.0.1
win-sqlserver: REFORM_SQL_INSTANCE ?= SQLEXPRESS
win-sqlserver: export REFORM_DATABASE = mssql
win-sqlserver: export REFORM_DRIVER = sqlserver
win-sqlserver: export REFORM_ROOT_SOURCE = sqlserver://$(REFORM_SQL_HOST)/$(REFORM_SQL_INSTANCE)
win-sqlserver: export REFORM_INIT_SOURCE = sqlserver://$(REFORM_SQL_HOST)/$(REFORM_SQL_INSTANCE)?database=reform-database
win-sqlserver: export REFORM_TEST_SOURCE = sqlserver://$(REFORM_SQL_HOST)/$(REFORM_SQL_INSTANCE)?database=reform-database
win-sqlserver: test
	mingw32-make test-db

.PHONY: docs parse reform reform-db
