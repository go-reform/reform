package bogus

//go:generate reform

// Bogus3 is used for testing. reform:bogus
type Bogus3 struct {
	bogus string `reform:"bogus"` // non-exported field with "reform:" tag should generate error
}
