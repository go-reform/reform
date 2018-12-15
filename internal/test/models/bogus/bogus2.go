package bogus

//go:generate reform

// Bogus2 is used for testing. reform:bogus
type Bogus2 struct {
	//nolint
	bogusType `reform:"bogus"` // anonymous field with "reform:" tag should generate error
}
