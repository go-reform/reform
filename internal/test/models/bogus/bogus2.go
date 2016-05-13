package bogus

//go:generate reform

// reform:bogus
type Bogus2 struct {
	bogusType `reform:"bogus"` // anonymous field with "reform:" tag should generate error
}
