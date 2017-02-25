package main

import (
	"fmt"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

// goTypeMySQL converts given SQL type to Go type. https://dev.mysql.com/doc/refman/5.7/en/data-types.html
func goTypeMySQL(sqlType string, nullable bool) (typ string, pack string, comment string) {
	switch sqlType {
	case "tinyint":
		return maybePointer("int8", nullable), "", ""
	case "smallint":
		return maybePointer("int16", nullable), "", ""
	case "mediumint", "int":
		return maybePointer("int32", nullable), "", ""
	case "bigint":
		return maybePointer("int64", nullable), "", ""

	case "decimal":
		return maybePointer("string", nullable), "", ""

	case "float":
		return maybePointer("float32", nullable), "", ""
	case "double":
		return maybePointer("float64", nullable), "", ""

	case "year", "date", "time", "datetime", "timestamp":
		return maybePointer("time.Time", nullable), "time", ""

	case "char", "varchar":
		fallthrough
	case "tinytext", "mediumtext", "text", "longtext":
		return maybePointer("string", nullable), "", ""

	case "binary", "varbinary":
		fallthrough
	case "tinyblob", "mediumblob", "blob", "longblob":
		return "[]byte", "", "" // never a pointer

	case "bool":
		return maybePointer("bool", nullable), "", ""

	default:
		// logger.Fatalf("unhandled MySQL type %q", sqlType)
		return "[]byte", "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}

func initModelsMySQL(db *reform.DB) (structs []StructData) {
	tables, err := db.SelectAllFrom(tableView, `WHERE table_schema = DATABASE()`)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	imports := make(map[string]struct{})
	for _, t := range tables {
		table := t.(*table)
		str := parse.StructInfo{Type: toCamelCase(table.TableName), SQLName: table.TableName}
		var comments []string

		key := getPrimaryKeyColumn(db, table.TableCatalog, table.TableSchema, table.TableName)

		tail := `WHERE table_catalog = ? AND table_schema = ? AND table_name = ? ORDER BY ordinal_position`
		columns, err := db.SelectAllFrom(columnView, tail, table.TableCatalog, table.TableSchema, table.TableName)
		if err != nil {
			logger.Fatalf("%s", err)
		}
		for i, c := range columns {
			column := c.(*column)
			typ, pack, comment := goTypeMySQL(column.Type, bool(column.IsNullable))
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
