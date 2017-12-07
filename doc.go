// Package reform is a better ORM for Go, based on non-empty interfaces and code generation.
//
// See README (https://github.com/go-reform/reform/blob/v1-stable/README.md) for quickstart information.
//
// Tagging
//
// reform allows one to add tags (comments) to generated queries with WithTag Querier method.
// They can be used to track queries from RDBMS logs and tools back to application code. For example, this code:
//  id := "baron"
//  person, err := DB.WithTag("GetProject:%v", id).FindByPrimaryKeyFrom(ProjectTable, id)
// will generate the following query:
//  SELECT /* GetProject:baron */ "projects"."name", "projects"."id", "projects"."start", "projects"."end" FROM "projects" WHERE "projects"."id" = ? LIMIT 1
// Please keep in mind that dynamic tags can affect RDBMS query cache. Consult your RDBMS documentation for details.
// Some known links:
//  MySQL / Percona Server: https://www.percona.com/doc/percona-server/LATEST/performance/query_cache_enhance.html#ignoring-comments
//  Microsoft SQL Server: https://msdn.microsoft.com/en-us/library/cc293623.aspx
//
// Short example
//
// This example shows some reform features.
// It uses https://github.com/AlekSi/pointer to get pointers to values of build-in types.
package reform // import "gopkg.in/reform.v1"

// Version defines reform version.
const Version = "v1.3.1"
