package main

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

// goTypeSQLite3 converts given SQL type to Go type. https://www.sqlite.org/datatype3.html
func goTypeSQLite3(sqlType string, nullable bool) (typ string, pack string, comment string) {
	// TODO is it a good logic? clarify, document it

	if strings.Contains(sqlType, "int") {
		return maybePointer("int64", nullable), "", ""
	}

	switch sqlType {
	case "character", "varchar", "varying character", "nchar", "native character", "nvarchar", "text", "clob":
		return maybePointer("string", nullable), "", ""
	case "blob", "":
		return "[]byte", "", "" // never a pointer
	case "real", "double", "double precision", "float":
		return maybePointer("float64", nullable), "", ""
	case "numeric", "decimal":
		return maybePointer("string", nullable), "", ""
	case "boolean":
		return maybePointer("bool", nullable), "", ""
	case "date", "datetime":
		return maybePointer("time.Time", nullable), "time", ""
	default:
		// logger.Fatalf("unhandled SQLite3 type %q", sqlType)
		return "[]byte", "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}

func initModelsSQLite3(db *reform.DB) (structs []StructData) {
	tables, err := db.SelectAllFrom(sqliteMasterView, "WHERE type = ?", "table")
	if err != nil {
		logger.Fatalf("%s", err)
	}

	imports := make(map[string]struct{})
	for _, table := range tables {
		tableName := table.(*sqliteMaster).Name
		if tableName == "sqlite_sequence" {
			continue
		}

		str := parse.StructInfo{Type: toCamelCase(tableName), SQLName: tableName}
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
				Name:   toCamelCase(column.Name),
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
