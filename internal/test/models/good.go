package models

import (
	"time"

	"github.com/AlekSi/pointer"

	"gopkg.in/reform.v1"
)

//go:generate reform

type (
	//reform:people
	Person struct {
		ID        int32      `reform:"id,pk"`
		GroupID   *int32     `reform:"group_id"`
		Name      string     `reform:"name"`
		Email     *string    `reform:"email"`
		CreatedAt time.Time  `reform:"created_at"`
		UpdatedAt *time.Time `reform:"updated_at"`
	}
)

// BeforeInsert sets CreatedAt if it's not set,
// then converts to UTC, truncates to second and strips monotonic clock reading from both CreatedAt and UpdatedAt.
func (p *Person) BeforeInsert() error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	p.CreatedAt = p.CreatedAt.UTC().Truncate(time.Second).AddDate(0, 0, 0)
	if p.UpdatedAt != nil {
		p.UpdatedAt = pointer.ToTime(p.UpdatedAt.UTC().Truncate(time.Second).AddDate(0, 0, 0))
	}

	return nil
}

// BeforeUpdate sets CreatedAt if it's not set,
// sets UpdatedAt,
// then converts to UTC, truncates to second and strips monotonic clock reading from both CreatedAt and UpdatedAt.
func (p *Person) BeforeUpdate() error {
	now := time.Now()

	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}

	p.UpdatedAt = &now

	p.CreatedAt = p.CreatedAt.UTC().Truncate(time.Second).AddDate(0, 0, 0)
	p.UpdatedAt = pointer.ToTime(p.UpdatedAt.UTC().Truncate(time.Second).AddDate(0, 0, 0))

	return nil
}

// AfterFind converts to UTC and truncates to second both CreatedAt and UpdatedAt.
func (p *Person) AfterFind() error {
	p.CreatedAt = p.CreatedAt.UTC().Truncate(time.Second)
	if p.UpdatedAt != nil {
		p.UpdatedAt = pointer.ToTime(p.UpdatedAt.UTC().Truncate(time.Second))
	}
	return nil
}

// Project represents row in table projects
// (reform:projects).
type Project struct {
	Name  string     `reform:"name"`
	ID    string     `reform:"id,pk"`
	Start time.Time  `reform:"start"`
	End   *time.Time `reform:"end"`
}

// BeforeInsert converts to UTC, truncates to day and strips monotonic clock reading from both Start and End.
func (p *Project) BeforeInsert() error {
	p.Start = p.Start.UTC().Truncate(24*time.Hour).AddDate(0, 0, 0)
	if p.End != nil {
		p.End = pointer.ToTime(p.End.UTC().Truncate(24*time.Hour).AddDate(0, 0, 0))
	}
	return nil
}

// BeforeUpdate converts to UTC, truncates to day and strips monotonic clock reading from both Start and End.
func (p *Project) BeforeUpdate() error {
	p.Start = p.Start.UTC().Truncate(24*time.Hour).AddDate(0, 0, 0)
	if p.End != nil {
		p.End = pointer.ToTime(p.End.UTC().Truncate(24*time.Hour).AddDate(0, 0, 0))
	}
	return nil
}

// AfterFind converts to UTC both Start and End.
func (p *Project) AfterFind() error {
	p.Start = p.Start.UTC()
	if p.End != nil {
		p.End = pointer.ToTime(p.End.UTC())
	}
	return nil
}

// PersonProject represents row in table person_project. reform:person_project
type PersonProject struct {
	PersonID  int32  `reform:"person_id"`
	ProjectID string `reform:"project_id"`
}

//reform:legacy.people
type LegacyPerson struct {
	ID   int32   `reform:"id,pk"`
	Name *string `reform:"name"`
}

// reform:id_only
type IDOnly struct {
	ID int32 `reform:"id,pk"`
}

// check interfaces
var (
	_ reform.BeforeInserter = (*Person)(nil)
	_ reform.BeforeUpdater  = (*Person)(nil)
	_ reform.AfterFinder    = (*Person)(nil)
	_ reform.BeforeInserter = (*Project)(nil)
	_ reform.BeforeUpdater  = (*Project)(nil)
	_ reform.AfterFinder    = (*Project)(nil)
)
