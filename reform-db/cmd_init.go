package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/parse"
)

func maybePointer(typ string, nullable bool) string {
	if nullable {
		return "*" + typ
	}
	return typ
}

func goType(sqlType string, nullable bool, dialect reform.Dialect) (string, string, string) {
	// https://www.postgresql.org/docs/current/static/datatype.html
	// https://dev.mysql.com/doc/refman/5.7/en/data-types.html
	// https://www.sqlite.org/datatype3.html
	// https://msdn.microsoft.com/en-us/library/ms187752.aspx

	// handle integer types
	switch dialect {
	case sqlite3.Dialect:
		switch sqlType {
		case "integer":
			return maybePointer("int64", nullable), "", ""
		}

	default:
		switch sqlType {
		case "tinyint":
			return maybePointer("int8", nullable), "", ""
		case "smallint":
			return maybePointer("int16", nullable), "", ""
		case "mediumint", "int", "integer":
			return maybePointer("int32", nullable), "", ""
		case "bigint":
			return maybePointer("int64", nullable), "", ""
		}
	}

	// order: PostgreSQL, MySQL, SQLite3, MS SQL
	switch sqlType {
	case "character", "character varying", "text":
		fallthrough
	case "char", "varchar", "tinytext", "mediumtext", "longtext":
		fallthrough
	case "nchar", "nvarchar":
		return maybePointer("string", nullable), "", ""

	// TODO blobs to []byte

	case "date", "time", "time with time zone", "timestamp", "timestamp with time zone":
		fallthrough
	case "datetime":
		return maybePointer("time.Time", nullable), "time", ""

	default:
		// never a pointer
		return "[]byte", "", fmt.Sprintf("// FIXME unhandled database type %q", sqlType)
	}
}

func toCamelCase(sqlName string) string {
	t := strings.Title(strings.Replace(sqlName, "_", " ", -1))
	return strings.Replace(t, " ", "", -1)
}

func getPrimaryKeyColumn(db *reform.DB, catalog, schema, name string) *keyColumnUsage {
	using := []string{"table_catalog", "table_schema", "table_name"}
	if db.Dialect == mysql.Dialect {
		// MySQL doesn't have table_catalog in table_constraints
		using = using[1:]
	}
	for i, u := range using {
		using[i] = fmt.Sprintf("key_column_usage.%s = table_constraints.%s", u, u)
	}
	q := fmt.Sprintf(`
		SELECT column_name, ordinal_position FROM information_schema.key_column_usage
			INNER JOIN information_schema.table_constraints ON %s
			WHERE key_column_usage.table_catalog = %s AND
				key_column_usage.table_schema = %s AND
				key_column_usage.table_name = %s AND
				constraint_type = 'PRIMARY KEY'
			ORDER BY ordinal_position DESC
		`, strings.Join(using, " AND "), db.Placeholder(1), db.Placeholder(2), db.Placeholder(3))
	row := db.QueryRow(q, catalog, schema, name)
	var key keyColumnUsage
	err := row.Scan(key.Pointers()...)
	if err == reform.ErrNoRows {
		err = nil
	}
	if err != nil {
		logger.Fatalf("%s", err)
	}
	if key.OrdinalPosition > 1 {
		logger.Fatalf("Expected single column primary key, got %d", key.OrdinalPosition)
	}
	return &key
}

func initModelsPostgreSQL(db *reform.DB) (structs []StructData) {
	tables, err := db.SelectAllFrom(tableView, `WHERE table_schema NOT IN ($1, $2)`, "pg_catalog", "information_schema")
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
			typ, pack, comment := goType(column.Type, bool(column.IsNullable), db.Dialect)
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
			typ, pack, comment := goType(column.Type, bool(column.IsNullable), db.Dialect)
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
			typ, pack, comment := goType(column.Type, !column.NotNull, db.Dialect)
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
			typ, pack, comment := goType(column.Type, bool(column.IsNullable), db.Dialect)
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

// cmdInit implements init command.
func cmdInit(db *reform.DB, dir string) {
	var structs []StructData
	switch db.Dialect {
	case postgresql.Dialect:
		structs = initModelsPostgreSQL(db)
	case mysql.Dialect:
		structs = initModelsMySQL(db)
	case sqlite3.Dialect:
		structs = initModelsSQLite3(db)
	case mssql.Dialect:
		structs = initModelsMSSQL(db)
	default:
		logger.Fatalf("unhandled dialect %s", db.Dialect)
	}

	pack := filepath.Base(dir)
	for _, s := range structs {
		logger.Debugf("%#v", s)

		f, err := os.Create(filepath.Join(dir, strings.ToLower(s.SQLName)+".go"))
		if err != nil {
			logger.Fatalf("%s", err)
		}

		logger.Debugf("Writing %s ...", f.Name())
		if _, err = f.WriteString("package " + pack + "\n"); err != nil {
			logger.Fatalf("%s", err)
		}
		if err = prologTemplate.Execute(f, s); err != nil {
			logger.Fatalf("%s", err)
		}
		if err = structTemplate.Execute(f, s); err != nil {
			logger.Fatalf("%s", err)
		}

		if err = f.Close(); err != nil {
			logger.Fatalf("%s", err)
		}
	}
}
