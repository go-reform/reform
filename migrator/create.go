package migrator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
)

// Create prepares an empty migration file with the given name and time in the given directory.
// If it is possible to create a migration (name and directory are valid and time is unique),
// the full path of the migration file will be returned.
func Create(dir, migration string, t time.Time) (string, error) {
	err := checkDir(dir)
	if err != nil {
		return "", err
	}
	err = checkName(migration)
	if err != nil {
		return "", err
	}

	dateTime := t.UTC().Format("2006-01-02-15-04")
	migrations, err := ioutil.ReadDir(dir)
	exists := false
	for _, f := range migrations {
		if strings.HasPrefix(f.Name(), dateTime) {
			exists = true
			break
		}
	}
	if exists {
		return "", fmt.Errorf("migration with the given time label (%s) already exists", dateTime)
	}

	filename := fmt.Sprintf("%s-%s.sql", dateTime, migration)
	f, err := os.Create(filepath.Join(dir, filename))
	if err != nil {

	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatalf("couldn't close migration file %s: %s", filename, err)
		}
	}()

	tplText := `-- +migration {{ .DateTime }}

-- place your code here
`
	tpl, err := template.New("migration").Parse(tplText)
	if err != nil {

	}

	err = tpl.Execute(f, map[string]string{"DateTime": dateTime})
	if err != nil {

	}
	return f.Name(), nil
}

// checkDir checks if the given migration directory exists and available.
func checkDir(dir string) error {
	if dir == "" {
		return fmt.Errorf("migraiton directory is not set")
	}

	fi, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("directory '%s' doesn't exist", dir)
		}
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("'%s' is not a directory", dir)
	}
	return nil
}

// checkMigrationName checks if the given migration name is valid.
func checkName(migration string) error {
	if len(migration) < 3 {
		return fmt.Errorf("migration name must containt at least 3 symbols (%d given)", len(migration))
	}
	if len(migration) > 100 {
		return fmt.Errorf("migration name mustn't contain more than 100 symbols (%d given)", len(migration))
	}
	pattern := `[a-z0-9\-]+`
	matched, err := regexp.MatchString(pattern, migration)
	if err != nil {
		return fmt.Errorf("couldn't match regex: %v", err)
	}
	if !matched {
		return fmt.Errorf(
			"migration name is incorrect," +
				" valid symbols are latin alphabet letters (in lower case), dashes(-) and digits",
		)
	}
	return nil
}
