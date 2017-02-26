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

test:
	rm -f internal/test/models/*_reform.go
	go install -v gopkg.in/reform.v1/...
	go test $(REFORM_TEST_FLAGS) -coverprofile=parse.cover gopkg.in/reform.v1/parse
	go generate -v -x gopkg.in/reform.v1/internal/test/models
	go install -v gopkg.in/reform.v1/internal/test/models

test-db:
	cat internal/test/sql/$(DATABASE)_init.sql \
		internal/test/sql/data.sql \
		internal/test/sql/$(DATABASE)_data.sql \
		internal/test/sql/$(DATABASE)_set.sql \
		| reform-db -db-driver="$(REFORM_DRIVER)" -db-source="$(REFORM_INIT_SOURCE)"
	go test $(REFORM_TEST_FLAGS) -coverprofile=$(REFORM_DRIVER).cover

check:
	-gometalinter ./... --deadline=60s --severity=vet:error

drone:
	drone exec --repo.trusted .drone-local.yml

postgres: export DATABASE = postgres
postgres: export REFORM_DRIVER = postgres
postgres: export REFORM_INIT_SOURCE = postgres://localhost/reform-database?sslmode=disable&TimeZone=UTC
postgres: export REFORM_TEST_SOURCE = postgres://localhost/reform-database?sslmode=disable&TimeZone=America/New_York
postgres: test
	-dropdb reform-database
	createdb reform-database
	make test-db

mysql: export DATABASE = mysql
mysql: export REFORM_DRIVER = mysql
mysql: export REFORM_INIT_SOURCE = root@/reform-database?parseTime=true&time_zone='UTC'&sql_mode='ANSI'&multiStatements=true
mysql: export REFORM_TEST_SOURCE = root@/reform-database?parseTime=true&time_zone='America%2FNew_York'
mysql: test
	echo 'DROP DATABASE IF EXISTS `reform-database`;' | mysql -uroot
	echo 'CREATE DATABASE `reform-database`;' | mysql -uroot
	make test-db

sqlite3: export DATABASE = sqlite3
sqlite3: export REFORM_DRIVER = sqlite3
sqlite3: export REFORM_INIT_SOURCE = /tmp/reform-database.sqlite3
sqlite3: export REFORM_TEST_SOURCE =/tmp/reform-database.sqlite3
sqlite3: test
	rm -f /tmp/reform-database.sqlite3
	make test-db

# this target is configured for Windows
mssql: REFORM_SQL_HOST ?= 127.0.0.1
mssql: REFORM_SQL_INSTANCE ?= SQLEXPRESS
mssql: SQLCMD = sqlcmd -b -I -S "$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE)"
mssql: export DATABASE = mssql
mssql: export REFORM_DRIVER = mssql
mssql: export REFORM_INIT_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE);database=reform-database
mssql: export REFORM_TEST_SOURCE = server=$(REFORM_SQL_HOST)\$(REFORM_SQL_INSTANCE);database=reform-database
mssql: test
	-$(SQLCMD) -Q "DROP DATABASE [reform-database];"
	$(SQLCMD) -Q "CREATE DATABASE [reform-database];"
	mingw32-make test-db

.PHONY: parse reform
