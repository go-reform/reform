package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mssql"
	"gopkg.in/reform.v1/dialects/mysql"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"gopkg.in/reform.v1/dialects/sqlserver"
	"gopkg.in/reform.v1/parse"
)

var (
	initFlags = flag.NewFlagSet("init", flag.ExitOnError)
	gofmtF    = initFlags.Bool("gofmt", true, "Format with gofmt")
)

func init() {
	initFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "`init` generates Go model files for existing database schema.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [global flags] init [init flags] [directory]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Global flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nInit flags:\n")
		initFlags.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
It uses information_schema or similar RDBMS mechanism to inspect database
structure. For each table, it generates a single file with single struct type
definition with fields, types, and tags. Generated code then should be checked
and edited manually.
`)
	}
}

func gofmt(path string) {
	if *gofmtF {
		cmd := exec.Command("gofmt", "-s", "-w", path)
		logger.Debugf(strings.Join(cmd.Args, " "))
		b, err := cmd.CombinedOutput()
		if err != nil {
			logger.Fatalf("gofmt error: %s", err)
		}
		logger.Debugf("gofmt output: %s", b)
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

// convertName converts snake_case name of table or column to CamelCase name of type or field.
// It also handles "_id" to "ID" conversion as a typical special case.
func convertName(sqlName string) string {
	fields := strings.Fields(strings.Replace(sqlName, "_", " ", -1))
	res := make([]string, len(fields))
	for i, f := range fields {
		if f == "id" {
			res[i] = "ID"
		} else {
			res[i] = strings.Title(f)
		}
	}
	return strings.Join(res, "")
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

	for _, t := range tables {
		imports := make(map[string]struct{})
		table := t.(*table)
		str := parse.StructInfo{
			Type:         convertName(table.TableName),
			SQLName:      table.TableName,
			PKFieldIndex: -1,
		}
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
				Name:   convertName(column.Name),
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
		// catalog is a currently selected database (reform-database, postgres, template0, etc.)
		// schema is a PostgreSQL schema (public, pg_catalog, information_schema, etc.)
		structs = initModelsInformationSchema(db, `WHERE table_schema = current_schema()`, goTypePostgres)
	case mysql.Dialect:
		// catalog is always "def"
		// schema is a database name (reform-database, information_schema, performance_schema, mysql, sys, etc.)
		structs = initModelsInformationSchema(db, `WHERE table_schema = DATABASE()`, goTypeMySQL)
	case sqlite3.Dialect:
		// SQLite is special
		structs = initModelsSQLite3(db)
	case mssql.Dialect, sqlserver.Dialect:
		// catalog is a currently selected database (reform-database, master, etc.)
		// schema is MS SQL schema (dbo, guest, sys, information_schema, etc.)
		structs = initModelsInformationSchema(db, `WHERE table_schema = SCHEMA_NAME()`, goTypeMSSQL)
	default:
		logger.Fatalf("unhandled dialect %s", db.Dialect)
	}

	// detect package name by importing package or from directory name
	var packageName string
	pack, err := build.ImportDir(dir, 0)
	if err == nil {
		packageName = pack.Name
	} else {
		s := strings.Split(filepath.Base(dir), ".")[0]
		packageName = strings.Replace(s, "-", "_", -1)
	}

	for _, s := range structs {
		logger.Debugf("%#v", s)

		f, err := os.Create(filepath.Join(dir, strings.ToLower(s.SQLName)+".go"))
		if err != nil {
			logger.Fatalf("%s", err)
		}

		logger.Debugf("Writing %s ...", f.Name())
		if _, err = f.WriteString("package " + packageName + "\n"); err != nil {
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

	gofmt(dir)
}
