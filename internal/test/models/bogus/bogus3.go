package bogus

//go:generate reform

// reform:bogus
type Bogus3 struct {
	bogus string `reform:"bogus"` // non-exported field with "reform:" tag should generate error
}
