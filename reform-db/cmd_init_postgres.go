package main

import (
	"fmt"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

// goTypePostgres converts given SQL type to Go type. https://www.postgresql.org/docs/current/static/datatype.html
func goTypePostgres(sqlType string, nullable bool) (typ string, pack string, comment string) {
	switch sqlType {
	case "smallint", "smallserial":
		return maybePointer("int16", nullable), "", ""
	case "integer", "serial":
		return maybePointer("int32", nullable), "", ""
	case "bigint", "bigserial":
		return maybePointer("int64", nullable), "", ""

	case "real":
		return maybePointer("float32", nullable), "", ""
	case "double precision":
		return maybePointer("float64", nullable), "", ""

	case "decimal", "numeric", "money":
		return maybePointer("string", nullable), "", ""

	case "character varying", "varchar", "character", "char", "text":
		return maybePointer("string", nullable), "", ""

	case "bytea":
		return "[]byte", "", "" // never a pointer

	case "timestamp", "timestamp with time zone", "date", "time", "time with time zone":
		return maybePointer("time.Time", nullable), "time", ""
		// interval can't be mapped to time.Duration: https://github.com/lib/pq/issues/78

	case "boolean":
		return maybePointer("bool", nullable), "", ""

	default:
		// logger.Fatalf("unhandled PostgreSQL type %q", sqlType)
		return "[]byte", "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}

func initModelsPostgreSQL(db *reform.DB) (structs []StructData) {
	tables, err := db.SelectAllFrom(tableView, `WHERE table_schema = current_schema()`)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	imports := make(map[string]struct{})
	for _, t := range tables {
		table := t.(*table)
		str := parse.StructInfo{Type: toCamelCase(table.TableName), SQLName: table.TableName}
		var comments []string

		key := getPrimaryKeyColumn(db, table.TableCatalog, table.TableSchema, table.TableName)

		tail := `WHERE table_catalog = $1 AND table_schema = $2 AND table_name = $3 ORDER BY ordinal_position`
		columns, err := db.SelectAllFrom(columnView, tail, table.TableCatalog, table.TableSchema, table.TableName)
		if err != nil {
			logger.Fatalf("%s", err)
		}
		for i, c := range columns {
			column := c.(*column)
			typ, pack, comment := goTypePostgres(column.Type, bool(column.IsNullable))
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
