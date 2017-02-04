package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/parse"
)

func goType(sqlType string) string {
	switch sqlType {
	case "integer":
		return "int"
	case "varchar":
		return "string"
	case "date", "time", "datetime":
		return "time.Time"
	default:
		return fmt.Sprintf("interface{} /* FIXME unhandled database type %q, please change it */", sqlType)
	}
}

func toCamelCase(sqlName string) string {
	t := strings.Title(strings.Replace(sqlName, "_", " ", -1))
	return strings.Replace(t, " ", "", -1)
}

func initModelsSQLite3(db *reform.DB) {
	tables, err := db.SelectAllFrom(sqliteMasterView, "WHERE type = ?", "table")
	if err != nil {
		logger.Fatal(err)
	}
	for _, table := range tables {
		tableName := table.(*sqliteMaster).Name
		if tableName == "sqlite_sequence" {
			continue
		}

		// logger.Printf("%s:", tableName)
		info := parse.StructInfo{Type: toCamelCase(tableName), SQLName: tableName}
		rows, err := db.Query("PRAGMA table_info(" + tableName + ")") // no placeholders for PRAGMA
		if err != nil {
			logger.Fatal(err)
		}
		for {
			var column sqliteTableInfo
			err = db.NextRow(&column, rows)
			if err != nil {
				break
			}
			// logger.Println(column)
			info.Fields = append(info.Fields, parse.FieldInfo{
				Name:   toCamelCase(column.Name),
				PKType: goType(column.Type), // FIXME this is Type, not PKType (not only PK)
				Column: column.Name,
			})
		}
		if err != reform.ErrNoRows {
			logger.Fatal(err)
		}
		rows.Close()
		// logger.Printf("%+v", info)

		if err = structTemplate.Execute(os.Stdout, info); err != nil {
			logger.Fatal(err)
		}
	}
}

func cmdInit(db *reform.DB, dialect reform.Dialect) {
	switch dialect {
	case sqlite3.Dialect:
		initModelsSQLite3(db)
	default:
		logger.Fatalf("unhandled dialect %s", dialect)
	}
}
