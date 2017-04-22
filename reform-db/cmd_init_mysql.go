package main

import (
	"fmt"
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
