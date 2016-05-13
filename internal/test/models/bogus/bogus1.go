package bogus

//go:generate reform

type (
	BogusType string
	bogusType string
)

// reform:bogus
type Bogus1 struct {
	BogusType `reform:"bogus"` // anonymous field with "reform:" tag should generate error
}
