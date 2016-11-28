all: check postgres mysql sqlite3

init:
	go get -u -d github.com/lib/pq \
				github.com/jackc/pgx/stdlib \
				github.com/go-sql-driver/mysql \
				github.com/ziutek/mymysql/... \
				github.com/mattn/go-sqlite3 \
				github.com/denisenkom/go-mssqldb

	go get -u -d github.com/AlekSi/pointer \
				github.com/kisielk/errcheck \
				github.com/golang/lint/golint \
				github.com/stretchr/testify/... \
				github.com/enodata/faker \
				github.com/mattn/goveralls

	go install -v github.com/kisielk/errcheck \
				github.com/golang/lint/golint \
				github.com/mattn/goveralls

test:
	rm -f internal/test/models/*_reform.go
	go install -v gopkg.in/reform.v1/...
	go test -coverprofile=parse.cover gopkg.in/reform.v1/parse
	go generate -v -x gopkg.in/reform.v1/internal/test/models
	go install -v gopkg.in/reform.v1/internal/test/models
	go test -i -v

check: test
	go vet ./...
	-errcheck ./...
	golint ./...

postgres: export REFORM_DRIVER = postgres
postgres: export REFORM_TEST_SOURCE = postgres://localhost/reform-test?sslmode=disable&TimeZone=America/New_York
postgres: test
	-dropdb reform-test
	createdb reform-test
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgres_init.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/data.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgres_data.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgres_set.sql
	go test -coverprofile=postgres.cover

# currently broken due to pgx's timezone handling
pgx: export REFORM_DRIVER = pgx
pgx: export REFORM_TEST_SOURCE = postgres://localhost/reform-test?sslmode=disable
pgx: test
	-dropdb reform-test
	createdb reform-test
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgres_init.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/data.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgres_data.sql
	env PGTZ=UTC psql -v ON_ERROR_STOP=1 -q -d reform-test < internal/test/sql/postgres_set.sql
	go test -coverprofile=pgx.cover

mysql: export REFORM_DRIVER = mysql
mysql: export REFORM_TEST_SOURCE = root@/reform-test?parseTime=true&strict=true&sql_notes=false&time_zone='America%2FNew_York'
mysql: test
	echo 'DROP DATABASE IF EXISTS `reform-test`;' | mysql -uroot
	echo 'CREATE DATABASE `reform-test`;' | mysql -uroot
	mysql -uroot reform-test < internal/test/sql/mysql_init.sql
	mysql -uroot reform-test < internal/test/sql/data.sql
	mysql -uroot reform-test < internal/test/sql/mysql_data.sql
	mysql -uroot reform-test < internal/test/sql/mysql_set.sql
	go test -coverprofile=mysql.cover

# currently broken due to mymysql's timezone handling
mymysql: export REFORM_DRIVER = mymysql
mymysql: export REFORM_TEST_SOURCE = reform-test/root/
mymysql: test
	echo 'DROP DATABASE IF EXISTS `reform-test`;' | mysql -uroot
	echo 'CREATE DATABASE `reform-test`;' | mysql -uroot
	mysql -uroot reform-test < internal/test/sql/mysql_init.sql
	mysql -uroot reform-test < internal/test/sql/data.sql
	mysql -uroot reform-test < internal/test/sql/mysql_data.sql
	mysql -uroot reform-test < internal/test/sql/mysql_set.sql
	go test -coverprofile=mymysql.cover

sqlite3: export REFORM_DRIVER = sqlite3
sqlite3: export REFORM_TEST_SOURCE = reform-test.sqlite3
sqlite3: test
	rm -f reform-test.sqlite3
	sqlite3 -bail reform-test.sqlite3 < internal/test/sql/sqlite3_init.sql
	sqlite3 -bail reform-test.sqlite3 < internal/test/sql/data.sql
	sqlite3 -bail reform-test.sqlite3 < internal/test/sql/sqlite3_data.sql
	sqlite3 -bail reform-test.sqlite3 < internal/test/sql/sqlite3_set.sql
	go test -coverprofile=sqlite3.cover

# this target is configured for Windows
mssql: REFORM_SQL_INSTANCE ?= 127.0.0.1\SQLEXPRESS
mssql: export REFORM_DRIVER = mssql
mssql: export REFORM_TEST_SOURCE = server=$(REFORM_SQL_INSTANCE);database=reform-test
mssql: test
	-sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -Q "DROP DATABASE [reform-test];"
	sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -Q "CREATE DATABASE [reform-test];"
	sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -d "reform-test" -i internal/test/sql/mssql_init.sql
	sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -d "reform-test" -i internal/test/sql/mssql_data.sql
	sqlcmd -b -I -S "$(REFORM_SQL_INSTANCE)" -d "reform-test" -i internal/test/sql/mssql_set.sql
	go test -coverprofile=mssql.cover
