package models

import "time"

//go:generate reform

//reform:schema_migrations
type Migration struct {
	Version string    `reform:"version,pk"`
	State   string    `reform:"state"`
	RunAt   time.Time `reform:"run_at"`
}
