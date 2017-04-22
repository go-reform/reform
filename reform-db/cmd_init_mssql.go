package main

import (
	"fmt"
)

// goTypeMSSQL converts given SQL type to Go type. https://msdn.microsoft.com/en-us/library/ms187752.aspx
func goTypeMSSQL(sqlType string, nullable bool) (typ string, pack string, comment string) {
	switch sqlType {
	case "tinyint":
		return maybePointer("uint8", nullable), "", "" // unsigned
	case "smallint":
		return maybePointer("int16", nullable), "", ""
	case "int":
		return maybePointer("int32", nullable), "", ""
	case "bigint":
		return maybePointer("int64", nullable), "", ""

	case "decimal", "numeric":
		return maybePointer("string", nullable), "", ""

	case "real":
		return maybePointer("float32", nullable), "", ""
	case "float":
		return maybePointer("float64", nullable), "", ""

	case "date", "time", "datetime", "datetime2", "smalldatetime":
		return maybePointer("time.Time", nullable), "time", ""

	case "char", "varchar", "text":
		fallthrough
	case "nchar", "nvarchar", "ntext":
		return maybePointer("string", nullable), "", ""

	case "binary", "varbinary":
		return "[]byte", "", "" // never a pointer

	case "bit":
		return maybePointer("bool", nullable), "", ""

	default:
		// logger.Fatalf("unhandled MSSQL type %q", sqlType)
		return "[]byte", "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}
