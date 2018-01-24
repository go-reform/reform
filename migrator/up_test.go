package migrator

import "testing"

func TestReadFromFiles(t *testing.T) {
	ReadFromFiles([]string{"../internal/test/migrations/*.sql"})
}
