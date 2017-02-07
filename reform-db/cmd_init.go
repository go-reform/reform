package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/parse"
)

func goType(sqlType string, dialect reform.Dialect) string {
	// https://www.postgresql.org/docs/current/static/datatype.html
	// https://dev.mysql.com/doc/refman/5.7/en/data-types.html
	// https://www.sqlite.org/datatype3.html

	// handle integer types
	switch dialect {
	case sqlite3.Dialect:
		switch sqlType {
		case "integer":
			return "int64"
		}

	default:
		switch sqlType {
		case "tinyint":
			return "int8"
		case "smallint":
			return "int16"
		case "mediumint", "int", "integer":
			return "int32"
		case "bigint":
			return "int64"
		}
	}

	// order: PostgreSQL, MySQL, SQLite3, MS SQL
	switch sqlType {
	case "character", "character varying", "text":
		fallthrough
	case "char", "varchar", "tinytext", "mediumtext", "longtext":
		fallthrough
	case "nchar", "nvarchar":
		return "string"

	// TODO blobs to []byte

	case "date", "time", "time with time zone", "timestamp", "timestamp with time zone":
		fallthrough
	case "datetime":
		return "time.Time"

	default:
		return fmt.Sprintf("interface{} /* FIXME unhandled database type %q, please change it */", sqlType)
	}
}

func toCamelCase(sqlName string) string {
	t := strings.Title(strings.Replace(sqlName, "_", " ", -1))
	return strings.Replace(t, " ", "", -1)
}

func initModelsPostgreSQL(db *reform.DB) (structs []parse.StructInfo) {
	tables, err := db.SelectAllFrom(tableView, `WHERE table_schema NOT IN ($1, $2)`, "pg_catalog", "information_schema")
	if err != nil {
		logger.Fatalf("%s", err)
	}

	for _, t := range tables {
		table := t.(*table)
		str := parse.StructInfo{Type: toCamelCase(table.Name), SQLName: table.Name}

		tail := `WHERE table_catalog = $1 AND table_schema = $2 AND table_name = $3 ORDER BY ordinal_position`
		columns, err := db.SelectAllFrom(columnView, tail, table.Catalog, table.Schema, table.Name)
		if err != nil {
			logger.Fatalf("%s", err)
		}
		for _, c := range columns {
			column := c.(*column)
			str.Fields = append(str.Fields, parse.FieldInfo{
				Name:   toCamelCase(column.Name),
				PKType: goType(column.Type, postgresql.Dialect), // FIXME this is Type, not PKType (not only PK)
				Column: column.Name,
			})
		}

		structs = append(structs, str)
	}

	return
}

func initModelsMySQL(db *reform.DB) (structs []parse.StructInfo) {
	tables, err := db.SelectAllFrom(tableView, `WHERE table_schema = DATABASE()`)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	for _, t := range tables {
		table := t.(*table)
		str := parse.StructInfo{Type: toCamelCase(table.Name), SQLName: table.Name}

		tail := `WHERE table_catalog = ? AND table_schema = ? AND table_name = ? ORDER BY ordinal_position`
		columns, err := db.SelectAllFrom(columnView, tail, table.Catalog, table.Schema, table.Name)
		if err != nil {
			logger.Fatalf("%s", err)
		}
		for _, c := range columns {
			column := c.(*column)
			str.Fields = append(str.Fields, parse.FieldInfo{
				Name:   toCamelCase(column.Name),
				PKType: goType(column.Type, mysql.Dialect), // FIXME this is Type, not PKType (not only PK)
				Column: column.Name,
			})
		}

		structs = append(structs, str)
	}

	return
}

func initModelsSQLite3(db *reform.DB) (structs []parse.StructInfo) {
	tables, err := db.SelectAllFrom(sqliteMasterView, "WHERE type = ?", "table")
	if err != nil {
		logger.Fatalf("%s", err)
	}

	for _, table := range tables {
		tableName := table.(*sqliteMaster).Name
		if tableName == "sqlite_sequence" {
			continue
		}

		str := parse.StructInfo{Type: toCamelCase(tableName), SQLName: tableName}
		rows, err := db.Query("PRAGMA table_info(" + tableName + ")") // no placeholders for PRAGMA
		if err != nil {
			logger.Fatalf("%s", err)
		}
		for {
			var column sqliteTableInfo
			if err = db.NextRow(&column, rows); err != nil {
				break
			}
			str.Fields = append(str.Fields, parse.FieldInfo{
				Name:   toCamelCase(column.Name),
				PKType: goType(column.Type, sqlite3.Dialect), // FIXME this is Type, not PKType (not only PK)
				Column: column.Name,
			})
		}
		if err != reform.ErrNoRows {
			logger.Fatalf("%s", err)
		}
		if err = rows.Close(); err != nil {
			logger.Fatalf("%s", err)
		}

		structs = append(structs, str)
	}

	return
}

func initModelsMSSQL(db *reform.DB) (structs []parse.StructInfo) {
	tables, err := db.SelectAllFrom(tableView, ``)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	for _, t := range tables {
		table := t.(*table)
		str := parse.StructInfo{Type: toCamelCase(table.Name), SQLName: table.Name}

		tail := `WHERE table_catalog = ? AND table_schema = ? AND table_name = ? ORDER BY ordinal_position`
		columns, err := db.SelectAllFrom(columnView, tail, table.Catalog, table.Schema, table.Name)
		if err != nil {
			logger.Fatalf("%s", err)
		}
		for _, c := range columns {
			column := c.(*column)
			str.Fields = append(str.Fields, parse.FieldInfo{
				Name:   toCamelCase(column.Name),
				PKType: goType(column.Type, mysql.Dialect), // FIXME this is Type, not PKType (not only PK)
				Column: column.Name,
			})
		}

		structs = append(structs, str)
	}

	return
}

func cmdInit(db *reform.DB, dialect reform.Dialect) {
	var structs []parse.StructInfo
	switch dialect {
	case postgresql.Dialect:
		structs = initModelsPostgreSQL(db)
	case mysql.Dialect:
		structs = initModelsMySQL(db)
	case sqlite3.Dialect:
		structs = initModelsSQLite3(db)
	case mssql.Dialect:
		structs = initModelsMSSQL(db)
	default:
		logger.Fatalf("unhandled dialect %s", dialect)
	}

	for _, s := range structs {
		if err := structTemplate.Execute(os.Stdout, s); err != nil {
			logger.Fatalf("%s", err)
		}
	}
}
