package main

import (
	"fmt"
)

// goTypeMySQL converts given SQL type to Go type. https://dev.mysql.com/doc/refman/5.7/en/data-types.html
func goTypeMySQL(sqlType string, nullable bool) (typ string, pack string, comment string) {
	switch sqlType {
	case sqlTypeTinyint:
		return maybePointer(typeInt8, nullable), "", ""
	case sqlTypeSmallint:
		return maybePointer(typeInt16, nullable), "", ""
	case sqlTypeMediumint, sqlTypeInt:
		return maybePointer(typeInt32, nullable), "", ""
	case sqlTypeBigint:
		return maybePointer(typeInt64, nullable), "", ""

	case sqlTypeDecimal:
		return maybePointer(typeString, nullable), "", ""

	case sqlTypeFloat:
		return maybePointer(typeFloat32, nullable), "", ""
	case sqlTypeDouble:
		return maybePointer(typeFloat64, nullable), "", ""

	case sqlTypeYear, sqlTypeDate, sqlTypeTime, sqlTypeDatetime, sqlTypeTimestamp:
		return maybePointer(typeTime, nullable), packageTime, ""

	case sqlTypeChar, sqlTypeVarchar:
		fallthrough
	case sqlTypeTinytext, sqlTypeMediumtext, sqlTypeText, sqlTypeLongtext:
		return maybePointer(typeString, nullable), "", ""

	case sqlTypeBinary, sqlTypeVarBinary:
		fallthrough
	case sqlTypeTinyblob, sqlTypeMediumblob, sqlTypeBlob, sqlTypeLongblob:
		return typeSliceByte, "", "" // never a pointer

	case sqlTypeBool:
		return maybePointer(typeBool, nullable), "", ""

	default:
		// logger.Fatalf("unhandled MySQL type %q", sqlType)
		return typeSliceByte, "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}
