package migrator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (m *Manager) Up(queries []string) error {
	err := m.createMigrationTable()
	if err != nil {
		return fmt.Errorf("couldn't prepare migration table: %v", err)
	}

	return nil
}

// ReadFromFile reads migrations from the given file paths. Supports templates, e.g. ./migrations/*.sql.
// Validation rules:
// - check if prefix of each file is a valid date
// - check if migration comment contains the same date as file prefix
func ReadFromFiles(files []string) ([]string, error) {
	var paths []string
	for _, f := range files {
		matches, err := filepath.Glob(f)
		if err != nil {
			return nil, fmt.Errorf("couldn't find matches for %s: %s ", f, err)
		}
		paths = append(paths, matches...)
	}

	for _, p := range paths {
		base := filepath.Base(p)
		if len(base) < 16 {
			fmt.Errorf("migration filename mustn't be less than 16 symbols, %d (%s) given", len(base), base)
		}
		prefix := base[:16]
		_, err := time.Parse("2006-01-02-15-04", prefix)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse date and time from migration name %s: %v", p, err)
		}

		f, err := os.Open(p)
		if err != nil {
			return nil, fmt.Errorf("couldn't open %s: %v", p, err)
		}
		line, _, err := bufio.NewReader(f).ReadLine()
		if err != nil {
			return nil, fmt.Errorf("couldn't read a line from %s: %v", p, err)
		}

		date, err := migrationDateFromComment(string(line))
		if err != nil {
			return nil, err
		}

		if date != prefix {
			return nil, fmt.Errorf(
				"the file prefix date (%s) and the file migration comment date (%s) are not matched: %s",
				p, prefix, date,
			)
		}
	}

	return paths, nil
}

func migrationDateFromComment(comment string) (string, error) {
	// the comment should look like "-- +migration YYYY-MM-DD-HH-mm",
	// so we are interested in [14:30] symbols
	if len(comment) != 30 {
		return "", fmt.Errorf("the right migration comment must contain 30 symbols")
	}
	date := comment[14:30]
	_, err := time.Parse("2006-01-02-15-04", date)
	if err != nil {
		return "", fmt.Errorf("couldn't parse date and time from migration comment %s: %v", comment, err)
	}
	return date, nil
}
