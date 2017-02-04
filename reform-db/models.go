package main

//go:generate reform

//reform:sqlite_master
type sqliteMaster struct {
	Name string `reform:"name"`
}

// TODO This "dummy" table name is ugly. We should do better.
//reform:dummy
type sqliteTableInfo struct {
	CID          int     `reform:"cid"`
	Name         string  `reform:"name"`
	Type         string  `reform:"type"`
	NotNull      bool    `reform:"notnull"`
	DefaultValue *string `reform:"dflt_value"`
	PK           bool    `reform:"pk"`
}
