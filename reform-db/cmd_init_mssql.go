package main

import (
	"fmt"
)

// goTypeMSSQL converts given SQL type to Go type. https://msdn.microsoft.com/en-us/library/ms187752.aspx
func goTypeMSSQL(sqlType string, nullable bool) (typ string, pack string, comment string) {
	switch sqlType {
	case sqlTypeTinyint:
		return maybePointer(typeUInt8, nullable), "", "" // unsigned
	case sqlTypeSmallint:
		return maybePointer(typeInt16, nullable), "", ""
	case sqlTypeInt:
		return maybePointer(typeInt32, nullable), "", ""
	case sqlTypeBigint:
		return maybePointer(typeInt64, nullable), "", ""

	case sqlTypeDecimal, sqlTypeNumeric:
		return maybePointer(typeString, nullable), "", ""

	case sqlTypeReal:
		return maybePointer(typeFloat32, nullable), "", ""
	case sqlTypeFloat:
		return maybePointer(typeFloat64, nullable), "", ""

	case sqlTypeDate, sqlTypeTime, sqlTypeDatetime, sqlTypeDatetime2, sqlTypeSmalldatetime:
		return maybePointer(typeTime, nullable), packageTime, ""

	case sqlTypeChar, sqlTypeVarchar, sqlTypeText:
		fallthrough
	case sqlTypeNChar, sqlTypeNVarChar, sqlTypeNText:
		return maybePointer(typeString, nullable), "", ""

	case sqlTypeBinary, sqlTypeVarBinary:
		return typeSliceByte, "", "" // never a pointer

	case sqlTypeBit:
		return maybePointer(typeBool, nullable), "", ""

	default:
		// logger.Fatalf("unhandled MSSQL type %q", sqlType)
		return typeSliceByte, "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType) // never a pointer
	}
}
