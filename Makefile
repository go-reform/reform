all: test test_lib_pq test_mattn_go-sqlite3 test_go-sql-driver_mysql

init:
	go get -u github.com/lib/pq
	go get -u github.com/jackc/pgx/stdlib
	go get -u github.com/mattn/go-sqlite3
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/ziutek/mymysql/...
	go get -u github.com/denisenkom/go-mssqldb
	go get -u github.com/AlekSi/pointer
	go get -u github.com/kisielk/errcheck
	go get -u github.com/golang/lint/golint
	go get -u github.com/stretchr/testify/...
	go get -u github.com/enodata/faker
	go get -u github.com/mattn/goveralls

install:
	rm -f internal/test/models/*_reform.go
	go install -v gopkg.in/reform.v1/...

test: install
	go test -coverprofile=parse.cover gopkg.in/reform.v1/parse
	go generate -v -x gopkg.in/reform.v1/internal/test/models
	go install -v gopkg.in/reform.v1/internal/test/models
	go test -i -v

check: test
	go vet ./...
	-errcheck ./...
	golint ./...

test_lib_pq: export REFORM_TEST_DRIVER = postgres
test_lib_pq: export REFORM_TEST_SOURCE = postgres://localhost:5432/reform-test?sslmode=disable&TimeZone=America/New_York
test_lib_pq:
	-dropdb reform-test
	createdb reform-test
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgresql_init.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/data.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgresql_data.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgresql_set.sql
	go test -coverprofile=test_lib_pq.cover

# currently broken due to pgx's timezone handling
test_jackc_pgx: export REFORM_TEST_DRIVER = pgx
test_jackc_pgx: export REFORM_TEST_SOURCE = postgres://localhost:5432/reform-test?sslmode=disable
test_jackc_pgx:
	-dropdb reform-test
	createdb reform-test
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgresql_init.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/data.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgresql_set.sql
	go test

test_mattn_go-sqlite3: export REFORM_TEST_DRIVER = sqlite3
test_mattn_go-sqlite3: export REFORM_TEST_SOURCE = reform-test.sqlite3
test_mattn_go-sqlite3:
	rm -f reform-test.sqlite3
	sqlite3 -bail reform-test.sqlite3 < internal/test/sql/sqlite3_init.sql
	sqlite3 -bail reform-test.sqlite3 < internal/test/sql/data.sql
	sqlite3 -bail reform-test.sqlite3 < internal/test/sql/sqlite3_set.sql
	go test -coverprofile=test_mattn_go-sqlite3.cover

test_go-sql-driver_mysql: export REFORM_TEST_DRIVER = mysql
test_go-sql-driver_mysql: export REFORM_TEST_SOURCE = root@/reform-test?parseTime=true&strict=true&sql_notes=false&time_zone='America%2FNew_York'
test_go-sql-driver_mysql:
	echo 'DROP DATABASE IF EXISTS `reform-test`;' | mysql -uroot
	echo 'CREATE DATABASE `reform-test`;' | mysql -uroot
	mysql -uroot reform-test < internal/test/sql/mysql_init.sql
	mysql -uroot reform-test < internal/test/sql/data.sql
	mysql -uroot reform-test < internal/test/sql/mysql_set.sql
	go test -coverprofile=test_go-sql-driver_mysql.cover

# currently broken due to mymysql's timezone handling
test_ziutek_mymysql: export REFORM_TEST_DRIVER = mymysql
test_ziutek_mymysql: export REFORM_TEST_SOURCE = reform-test/root/
test_ziutek_mymysql:
	echo 'DROP DATABASE IF EXISTS `reform-test`;' | mysql -uroot
	echo 'CREATE DATABASE `reform-test`;' | mysql -uroot
	mysql -uroot reform-test < internal/test/sql/mysql_init.sql
	mysql -uroot reform-test < internal/test/sql/data.sql
	mysql -uroot reform-test < internal/test/sql/mysql_set.sql
	go test

# this target is configured for Windows
test_denisenkom_go-mssqldb: REFORM_SQL_INSTANCE ?= 127.0.0.1\SQLEXPRESS
test_denisenkom_go-mssqldb: export REFORM_TEST_DRIVER = mssql
test_denisenkom_go-mssqldb: export REFORM_TEST_SOURCE = server=$(REFORM_SQL_INSTANCE);database=reform-test
test_denisenkom_go-mssqldb:
	-sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -Q "DROP DATABASE [reform-test];"
	sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -Q "CREATE DATABASE [reform-test];"
	sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -d "reform-test" -i internal/test/sql/mssql_init.sql
	sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -d "reform-test" -i internal/test/sql/mssql_data.sql
	sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -d "reform-test" -i internal/test/sql/mssql_set.sql
	go test -coverprofile=test_denisenkom_go-mssqldb.cover

parse:
	# nothing, hack for our Travis-CI configuration
	# see 'test' target here and $TARGET in .travis.yml
