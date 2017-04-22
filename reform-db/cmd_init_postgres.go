package main

import (
	"fmt"
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

	case "decimal", "numeric":
		return maybePointer("string", nullable), "", ""

	case "real":
		return maybePointer("float32", nullable), "", ""
	case "double precision":
		return maybePointer("float64", nullable), "", ""

	case "character varying", "varchar", "character", "char", "text":
		return maybePointer("string", nullable), "", ""

	case "bytea":
		return "[]byte", "", "" // never a pointer

	case "date", "time", "time with time zone", "timestamp", "timestamp with time zone":
		return maybePointer("time.Time", nullable), "time", ""
		// interval can't be mapped to time.Duration: https://github.com/lib/pq/issues/78

	case "boolean":
		return maybePointer("bool", nullable), "", ""

	default:
		// logger.Fatalf("unhandled PostgreSQL type %q", sqlType)
		return "[]byte", "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}
