package main

import (
	"fmt"
)

// goTypePostgres converts given SQL type to Go type. https://www.postgresql.org/docs/current/static/datatype.html
func goTypePostgres(sqlType string, nullable bool) (typ string, pack string, comment string) {
	switch sqlType {
	case sqlTypeSmallint, sqlTypeSmallserial:
		return maybePointer(typeInt16, nullable), "", ""
	case sqlTypeInteger, sqlTypeSerial:
		return maybePointer(typeInt32, nullable), "", ""
	case sqlTypeBigint, sqlTypeBigserial:
		return maybePointer(typeInt64, nullable), "", ""

	case sqlTypeDecimal, sqlTypeNumeric:
		return maybePointer(typeString, nullable), "", ""

	case sqlTypeReal:
		return maybePointer(typeFloat32, nullable), "", ""
	case sqlTypeDoublePrecision:
		return maybePointer(typeFloat64, nullable), "", ""

	case sqlTypeCharacterVarying, sqlTypeVarchar, sqlTypeCharacter, sqlTypeChar, sqlTypeText:
		return maybePointer(typeString, nullable), "", ""

	case sqlTypeBytea:
		return typeSliceByte, "", "" // never a pointer

	case sqlTypeDate, sqlTypeTime, sqlTypeTimeWithTimeZone, sqlTypeTimestamp, sqlTypeTimestampWithTimeZone:
		return maybePointer(typeTime, nullable), packageTime, ""
		// interval can't be mapped to time.Duration: https://github.com/lib/pq/issues/78

	case sqlTypeBoolean:
		return maybePointer(typeBool, nullable), "", ""

	default:
		// logger.Fatalf("unhandled PostgreSQL type %q", sqlType)
		return typeSliceByte, "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}
