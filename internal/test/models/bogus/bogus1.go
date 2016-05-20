package bogus

//go:generate reform

type (
	// BogusType is an exported type used for testing.
	BogusType string

	// BogusType is a non-exported type used for testing.
	bogusType string
)

// Bogus1 is used for testing. reform:bogus
type Bogus1 struct {
	BogusType `reform:"bogus"` // anonymous field with "reform:" tag should generate error
}
