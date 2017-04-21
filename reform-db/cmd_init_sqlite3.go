package main

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

// goTypeSQLite3 converts given SQL type to Go type. https://www.sqlite.org/datatype3.html
func goTypeSQLite3(sqlType string, nullable bool) (typ string, pack string, comment string) {
	// SQLite3 has quite unique dynamic type system with storage classes and type affinities.
	// In short:
	// * table columns don't have rigid types;
	// * value has storage class (null, integer, real, text, blob), which defines how value is stored on disk;
	// * table column has type affinity (text, numeric, integer, real, blob), which defines preferred storage class
	//   for values in this column;
	// * table column declared type in CREATE TABLE defines column affinity with a set of rules;
	// * table_info returns column declared type;
	// * we try to mirror SQLite's set of rules of defining column affinity from declared type to define Go type;
	// * we also extend this set of rules with some common SQL data types;
	// * it's not 100% accurate (because SQLite is dynamically typed), but it follows actual and best practices.

	sqlType = strings.ToLower(sqlType)

	// SQLite rules 1-4
	switch {
	case strings.Contains(sqlType, "int"):
		return maybePointer("int64", nullable), "", ""

	case strings.Contains(sqlType, "char") || strings.Contains(sqlType, "clob") || strings.Contains(sqlType, "text"):
		return maybePointer("string", nullable), "", ""

	case strings.Contains(sqlType, "blob") || sqlType == "":
		return "[]byte", "", "" // never a pointer

	case strings.Contains(sqlType, "real") || strings.Contains(sqlType, "floa") || strings.Contains(sqlType, "doub"):
		return maybePointer("float64", nullable), "", ""
	}

	// common SQL data types
	switch {
	// numeric, decimal, etc.
	case strings.Contains(sqlType, "num") || strings.Contains(sqlType, "dec"):
		return maybePointer("string", nullable), "", ""

	// bool, boolean, etc.
	case strings.Contains(sqlType, "bool"):
		return maybePointer("bool", nullable), "", ""

	// date, datetime, timestamp, etc.
	case strings.Contains(sqlType, "date") || strings.Contains(sqlType, "time"):
		return maybePointer("time.Time", nullable), "time", ""

	default:
		// logger.Fatalf("unhandled SQLite3 type %q", sqlType)
		return "[]byte", "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}

// initModelsSQLite3 returns structs from SQLite3 database.
func initModelsSQLite3(db *reform.DB) (structs []StructData) {
	tables, err := db.SelectAllFrom(sqliteMasterView, "WHERE type = ?", "table")
	if err != nil {
		logger.Fatalf("%s", err)
	}

	for _, table := range tables {
		imports := make(map[string]struct{})
		tableName := table.(*sqliteMaster).Name
		if tableName == "sqlite_sequence" {
			continue
		}

		str := parse.StructInfo{
			Type:         convertName(tableName),
			SQLName:      tableName,
			PKFieldIndex: -1,
		}
		var comments []string

		rows, err := db.Query("PRAGMA table_info(" + tableName + ")") // no placeholders for PRAGMA
		if err != nil {
			logger.Fatalf("%s", err)
		}
		for {
			var column sqliteTableInfo
			if err = db.NextRow(&column, rows); err != nil {
				break
			}
			if column.PK {
				str.PKFieldIndex = len(str.Fields)
			}
			typ, pack, comment := goTypeSQLite3(column.Type, !column.NotNull)
			if pack != "" {
				imports[pack] = struct{}{}
			}
			comments = append(comments, comment)
			str.Fields = append(str.Fields, parse.FieldInfo{
				Name:   convertName(column.Name),
				Type:   typ,
				Column: column.Name,
			})
		}
		if err != reform.ErrNoRows {
			logger.Fatalf("%s", err)
		}
		if err = rows.Close(); err != nil {
			logger.Fatalf("%s", err)
		}

		structs = append(structs, StructData{
			Imports:       imports,
			StructInfo:    str,
			FieldComments: comments,
		})
	}

	return
}
