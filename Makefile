all: test postgres mysql sqlite3 check

# extra flags like -v
REFORM_TEST_FLAGS ?=

# SHELL = go run .github/shell.go

download_deps:
	# download drivers
	go get -v -u -d github.com/lib/pq \
		github.com/go-sql-driver/mysql \
		github.com/mattn/go-sqlite3 \
		github.com/denisenkom/go-mssqldb

	# download other deps
	go get -v -u -d github.com/AlekSi/pointer \
		github.com/stretchr/testify/... \
		github.com/enodata/faker \
		github.com/alecthomas/gometalinter \
		github.com/AlekSi/goveralls

	# download linters
	go install -v github.com/alecthomas/gometalinter
	gometalinter --install --update --download-only

install_deps:
	go install -v github.com/alecthomas/gometalinter \
		github.com/AlekSi/goveralls
	gometalinter --install
	go test -i -v

# run unit tests, generate models, install tools
test:
	rm -f *.cover
	rm -f internal/test/models/*_reform.go
	rm -f reform-db/*_reform.go

	go install -v gopkg.in/reform.v1/reform
	go test $(REFORM_TEST_FLAGS) -coverprofile=parse.cover gopkg.in/reform.v1/parse
	go generate -v -x gopkg.in/reform.v1/internal/test/models
	go install -v gopkg.in/reform.v1/internal/test/models

	go generate -v -x gopkg.in/reform.v1/reform-db
	go install -v gopkg.in/reform.v1/reform-db

# initialize database and run tests
test-db:
	-reform-db -db-driver="$(REFORM_DRIVER)" -db-source="$(REFORM_ROOT_SOURCE)" exec \
		internal/test/sql/$(DATABASE)_drop.sql
	reform-db -db-driver="$(REFORM_DRIVER)" -db-source="$(REFORM_ROOT_SOURCE)" exec \
		internal/test/sql/$(DATABASE)_create.sql
	reform-db -db-driver="$(REFORM_DRIVER)" -db-source="$(REFORM_INIT_SOURCE)" exec \
		internal/test/sql/$(DATABASE)_init.sql \
		internal/test/sql/data.sql \
		internal/test/sql/$(DATABASE)_data.sql \
		internal/test/sql/$(DATABASE)_set.sql
	go test $(REFORM_TEST_FLAGS) -coverprofile=$(REFORM_DRIVER)-reform-db.cover gopkg.in/reform.v1/reform-db
	go test $(REFORM_TEST_FLAGS) -coverprofile=$(REFORM_DRIVER).cover

check:
	-gometalinter ./... --deadline=180s --severity=vet:error

drone:
	drone exec --repo.trusted .drone-local.yml

# create local PostgreSQL database and run tests
postgres: export DATABASE = postgres
postgres: export REFORM_DRIVER = postgres
postgres: export REFORM_ROOT_SOURCE = postgres://localhost/template1?sslmode=disable
postgres: export REFORM_INIT_SOURCE = postgres://localhost/reform-database?sslmode=disable&TimeZone=UTC
postgres: export REFORM_TEST_SOURCE = postgres://localhost/reform-database?sslmode=disable&TimeZone=America/New_York
postgres: test
	make test-db

# create local MySQL database and run tests
mysql: export DATABASE = mysql
mysql: export REFORM_DRIVER = mysql
mysql: export REFORM_ROOT_SOURCE = root@/mysql
mysql: export REFORM_INIT_SOURCE = root@/reform-database?parseTime=true&time_zone='UTC'&sql_mode='ANSI'&multiStatements=true
mysql: export REFORM_TEST_SOURCE = root@/reform-database?parseTime=true&time_zone='America%2FNew_York'
mysql: test
	make test-db

# create local SQLite3 database and run tests
sqlite3: export DATABASE = sqlite3
sqlite3: export REFORM_DRIVER = sqlite3
sqlite3: export REFORM_ROOT_SOURCE = /tmp/reform-database.sqlite3
sqlite3: export REFORM_INIT_SOURCE = /tmp/reform-database.sqlite3
sqlite3: export REFORM_TEST_SOURCE = /tmp/reform-database.sqlite3
sqlite3: test
	rm -f /tmp/reform-database.sqlite3
	make test-db

# create SQL Server database and run tests with mssql driver (Windows only)
mssql: REFORM_SQL_HOST ?= 127.0.0.1
mssql: REFORM_SQL_INSTANCE ?= SQLEXPRESS
mssql: export DATABASE = mssql
mssql: export REFORM_DRIVER = mssql
mssql: export REFORM_ROOT_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE)
mssql: export REFORM_INIT_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE);database=reform-database
mssql: export REFORM_TEST_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE);database=reform-database
mssql: test
	mingw32-make test-db

# create SQL Server database and run tests with sqlserver driver (Windows only)
sqlserver: REFORM_SQL_HOST ?= 127.0.0.1
sqlserver: REFORM_SQL_INSTANCE ?= SQLEXPRESS
sqlserver: export DATABASE = mssql
sqlserver: export REFORM_DRIVER = sqlserver
sqlserver: export REFORM_ROOT_SOURCE = sqlserver://$(REFORM_SQL_HOST)/$(REFORM_SQL_INSTANCE)
sqlserver: export REFORM_INIT_SOURCE = sqlserver://$(REFORM_SQL_HOST)/$(REFORM_SQL_INSTANCE)?database=reform-database
sqlserver: export REFORM_TEST_SOURCE = sqlserver://$(REFORM_SQL_HOST)/$(REFORM_SQL_INSTANCE)?database=reform-database
sqlserver: test
	mingw32-make test-db

.PHONY: parse reform reform-db
