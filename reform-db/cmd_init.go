package main

import (
	"flag"
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

var (
	initFlags = flag.NewFlagSet("init", flag.ExitOnError)
)

func init() {
	initFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "`init` generates Go model files for existing database schema.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [global flags] init [directory]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Global flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "TODO.\n")
		initFlags.PrintDefaults()
	}
}

type typeFunc func(sqlType string, nullable bool) (typ string, pack string, comment string)

// maybePointer returns *typ if nullable, typ otherwise.
func maybePointer(typ string, nullable bool) string {
	if nullable {
		return "*" + typ
	}
	return typ
}

// toCamelCase converts snake_case string (for example, table or column name)
// to CamelCase (for example, type or field name).
func toCamelCase(sqlName string) string {
	t := strings.Title(strings.Replace(sqlName, "_", " ", -1))
	return strings.Replace(t, " ", "", -1)
}

// getPrimaryKeyColumn returns single primary key column for given table, or nil.
func getPrimaryKeyColumn(db *reform.DB, catalog, schema, tableName string) *keyColumnUsage {
	using := []string{"table_catalog", "table_schema", "table_name"}
	if db.Dialect == mysql.Dialect {
		// MySQL doesn't have table_catalog in table_constraints
		using = using[1:]
	}
	for i, u := range using {
		using[i] = fmt.Sprintf("key_column_usage.%s = table_constraints.%s", u, u)
	}
	q := fmt.Sprintf(
		`SELECT column_name, ordinal_position FROM information_schema.key_column_usage
			INNER JOIN information_schema.table_constraints ON %s
			WHERE key_column_usage.table_catalog = %s AND
				key_column_usage.table_schema = %s AND
				key_column_usage.table_name = %s AND
				constraint_type = 'PRIMARY KEY'
			ORDER BY ordinal_position DESC`,
		strings.Join(using, " AND "), db.Placeholder(1), db.Placeholder(2), db.Placeholder(3),
	)
	row := db.QueryRow(q, catalog, schema, tableName)
	var key keyColumnUsage
	if err := row.Scan(key.Pointers()...); err != nil {
		if err == reform.ErrNoRows {
			return nil
		}
		logger.Fatalf("%s", err)
	}
	if key.OrdinalPosition > 1 {
		logger.Fatalf("Expected single column primary key, got %d", key.OrdinalPosition)
	}
	return &key
}

// initModelsInformationSchema returns structs from database with information_schema.
func initModelsInformationSchema(db *reform.DB, tablesTail string, typeFunc typeFunc) (structs []StructData) {
	tables, err := db.SelectAllFrom(tableView, tablesTail)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	imports := make(map[string]struct{})
	for _, t := range tables {
		table := t.(*table)
		str := parse.StructInfo{Type: toCamelCase(table.TableName), SQLName: table.TableName}
		var comments []string

		key := getPrimaryKeyColumn(db, table.TableCatalog, table.TableSchema, table.TableName)

		tail := fmt.Sprintf(
			`WHERE table_catalog = %s AND table_schema = %s AND table_name = %s ORDER BY ordinal_position`,
			db.Placeholder(1), db.Placeholder(2), db.Placeholder(3),
		)
		columns, err := db.SelectAllFrom(columnView, tail, table.TableCatalog, table.TableSchema, table.TableName)
		if err != nil {
			logger.Fatalf("%s", err)
		}
		for i, c := range columns {
			column := c.(*column)
			typ, pack, comment := typeFunc(column.Type, bool(column.IsNullable))
			if pack != "" {
				imports[pack] = struct{}{}
			}
			comments = append(comments, comment)
			str.Fields = append(str.Fields, parse.FieldInfo{
				Name:   toCamelCase(column.Name),
				Type:   typ,
				Column: column.Name,
			})

			if key != nil && key.ColumnName == column.Name {
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
		structs = initModelsInformationSchema(db, `WHERE table_schema = current_schema()`, goTypePostgres)
	case mysql.Dialect:
		structs = initModelsInformationSchema(db, `WHERE table_schema = DATABASE()`, goTypeMySQL)
	case sqlite3.Dialect:
		structs = initModelsSQLite3(db)
	case mssql.Dialect:
		structs = initModelsInformationSchema(db, ``, goTypeMSSQL)
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
