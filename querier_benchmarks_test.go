// +build bench

package reform_test

import (
	"testing"
	"time"

	"github.com/enodata/faker"

	"gopkg.in/reform.v1"
	. "gopkg.in/reform.v1/internal/test/models"
)

var err error

const N = 10000

// FindByPrimaryKeyFrom

func benchmarkFindByPrimaryKeyFrom(tb testing.TB, tx *reform.TX, n int) {
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

// Update

func benchmarkUpdate(tb testing.TB, tx *reform.TX, n int, person *Person) {
	email1 := "email1"
	email2 := "email2"

	for i := 0; i < n; i++ {
		if i%2 == 0 {
			person.Email = &email1
		} else {
			person.Email = &email2
		}

		err = tx.Update(person)
		if err != nil {
			tb.Fatal(err)
		}
	}
}

func TestBenchmarkUpdate(t *testing.T) {
	tx := setupTransaction(t, false)
	defer tearDownTransaction(t, tx)

	person := &Person{ID: 1}
	err := tx.Reload(person)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	benchmarkUpdate(t, tx, N, person)
	t.Logf("N = %d in %s", N, time.Now().Sub(start))
}

func BenchmarkUpdate(b *testing.B) {
	tx := setupTransaction(b, false)
	defer tearDownTransaction(b, tx)

	person := &Person{ID: 1}
	err = tx.Reload(person)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	benchmarkUpdate(b, tx, b.N, person)
	b.StopTimer()
}
