// +build bench

package reform_test

import (
	"testing"
	"time"

	"github.com/enodata/faker"

	"gopkg.in/reform.v1"
	. "gopkg.in/reform.v1/internal/test/models"
)

// FindByPrimaryKeyFrom

func benchmarkFindByPrimaryKeyFrom(tb testing.TB, tx *reform.TX, n int) {
	var err error
	for i := 0; i < n; i++ {
		_, err = tx.FindByPrimaryKeyFrom(PersonTable, 1)
		if err != nil {
			tb.Fatal(err)
		}
	}
}

func TestBenchmarkFindByPrimaryKeyFrom(t *testing.T) {
	tx := setupTransaction(t, false)
	defer tearDownTransaction(t, tx)

	const N = 10000
	start := time.Now()
	benchmarkFindByPrimaryKeyFrom(t, tx, N)
	t.Logf("N = %d in %s", N, time.Now().Sub(start))
}

func BenchmarkFindByPrimaryKeyFrom(b *testing.B) {
	tx := setupTransaction(b, false)
	defer tearDownTransaction(b, tx)

	b.ResetTimer()
	benchmarkFindByPrimaryKeyFrom(b, tx, b.N)
	b.StopTimer()
}

// Insert

func benchmarkInsert(tb testing.TB, tx *reform.TX, n int, newEmail string) {
	var err error
	for i := 0; i < n; i++ {
		err = tx.Insert(&Person{Email: &newEmail})
		if err != nil {
			tb.Fatal(err)
		}
	}
}

func TestBenchmarkInsert(t *testing.T) {
	tx := setupTransaction(t, false)
	defer tearDownTransaction(t, tx)

	newEmail := faker.Internet().Email()

	const N = 10000
	start := time.Now()
	benchmarkInsert(t, tx, N, newEmail)
	t.Logf("N = %d in %s", N, time.Now().Sub(start))
}

func BenchmarkInsert(b *testing.B) {
	tx := setupTransaction(b, false)
	defer tearDownTransaction(b, tx)

	newEmail := faker.Internet().Email()

	b.ResetTimer()
	benchmarkInsert(b, tx, b.N, newEmail)
	b.StopTimer()
}
