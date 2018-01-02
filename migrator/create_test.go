package migrator

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	dir := os.TempDir()
	migration := "a-couple-of-words"
	filename := filepath.Join(dir, "2017-12-22-17-30-"+migration+".sql")

	// if file already exist, delete it to be able to run a test in the clean env
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		require.NoError(t, err)
	}

	dateTime := time.Date(2017, 12, 22, 17, 30, 50, 55, time.UTC)
	path, err := Create(dir, migration, dateTime)
	require.NoError(t, err)

	_, err = os.Stat(filename)
	require.NoError(t, err)

	expectedData := []byte(`-- +migration 2017-12-22-17-30

-- place your code here
`)
	data, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	res := bytes.Compare(expectedData, data)
	require.Equal(t, 0, res, "Generated migration file content doesn't match expected data")
	require.Equal(t, filename, path, "Generated migration path doesn't match expected path")

	// try to create a migration with the same date second time (an error is expected)
	migration2 := "a-couple-of-words-2"
	filename2 := filepath.Join(dir, "2017-12-22-17-30-"+migration2+".sql")
	if _, err := os.Stat(filename2); err == nil {
		err = os.Remove(filename2)
		require.NoError(t, err)
	}
	path2, err := Create(dir, migration2, dateTime)
	require.Error(t, err)
	require.Equal(t, "", path2)
}

func TestCheckDir(t *testing.T) {
	set := []struct {
		dir   string
		valid bool
	}{
		{"/does-not-exist-" + strconv.Itoa(rand.Int()), false}, // doesn't exists
		{os.TempDir(), true},
		{"", false}, // empty
	}
	for i, item := range set {
		err := checkDir(item.dir)
		if item.valid {
			require.NoError(t, err, "Case %d: %v", i, item)
		} else {
			require.Error(t, err, "Case %d: %v", i, item)
		}
	}
}

func TestCheckName(t *testing.T) {
	set := []struct {
		name  string
		valid bool
	}{
		{"abc", true},
		{"12abc-def-349234-zdfd6-", true},
		{"", false},   // empty
		{"ab", false}, // too short
		{func() string {
			var str string
			for i := 0; i <= 100; i++ {
				str += "a"
			}
			return str
		}(), false}, // too long
	}
	for i, item := range set {
		err := checkName(item.name)
		if item.valid {
			require.NoError(t, err, "Case %d: %v", i, item)
		} else {
			require.Error(t, err, "Case %d: %v", i, item)
		}
	}
}
