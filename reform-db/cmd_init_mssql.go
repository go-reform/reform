package main

import (
	"fmt"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

// goTypeMSSQL converts given SQL type to Go type. https://msdn.microsoft.com/en-us/library/ms187752.aspx
func goTypeMSSQL(sqlType string, nullable bool) (typ string, pack string, comment string) {
	// order: PostgreSQL, MySQL, SQLite3, MS SQL
	switch sqlType {
	case "tinyint":
		return maybePointer("uint8", nullable), "", "" // unsigned
	case "smallint":
		return maybePointer("int16", nullable), "", ""
	case "int":
		return maybePointer("int32", nullable), "", ""
	case "bigint":
		return maybePointer("int64", nullable), "", ""

		// TODO

	default:
		// logger.Fatalf("unhandled MSSQL type %q", sqlType)
		return "[]byte", "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}

func initModelsMSSQL(db *reform.DB) (structs []StructData) {
	tables, err := db.SelectAllFrom(tableView, ``)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	imports := make(map[string]struct{})
	for _, t := range tables {
		table := t.(*table)
		str := parse.StructInfo{Type: toCamelCase(table.TableName), SQLName: table.TableName}

		key := getPrimaryKeyColumn(db, table.TableCatalog, table.TableSchema, table.TableName)
		var comments []string

		tail := `WHERE table_catalog = ? AND table_schema = ? AND table_name = ? ORDER BY ordinal_position`
		columns, err := db.SelectAllFrom(columnView, tail, table.TableCatalog, table.TableSchema, table.TableName)
		if err != nil {
			logger.Fatalf("%s", err)
		}
		for i, c := range columns {
			column := c.(*column)
			typ, pack, comment := goTypeMSSQL(column.Type, bool(column.IsNullable))
			if pack != "" {
				imports[pack] = struct{}{}
			}
			comments = append(comments, comment)
			str.Fields = append(str.Fields, parse.FieldInfo{
				Name:   toCamelCase(column.Name),
				Type:   typ,
				Column: column.Name,
			})

			if key.ColumnName == column.Name {
				str.PKFieldIndex = i
			}
		}

		structs = append(structs, StructData{
			Imports:       imports,
			StructInfo:    str,
			FieldComments: comments,
		})
	}

	return
}
