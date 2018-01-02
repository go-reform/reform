package migrator

import "fmt"

func (m *Manager) Up(queries []string) error {
	err := m.createMigrationTable()
	if err != nil {
		return fmt.Errorf("couldn't prepare migration table: %v", err)
	}

	return nil
}
