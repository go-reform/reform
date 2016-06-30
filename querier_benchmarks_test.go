package reform_test

import (
	"testing"

	"github.com/enodata/faker"

	. "gopkg.in/reform.v1/internal/test/models"
)

func BenchmarkInsert(b *testing.B) {
	tx := setupTransaction(b, false)
	defer tearDownTransaction(b, tx)

	newEmail := faker.Internet().Email()
	var err error

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tx.Insert(&Person{Email: &newEmail})
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}
