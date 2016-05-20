package bogus

//go:generate reform

// Bogus7 is used for testing. reform:bogus
type Bogus7 struct {
	Bogus *string `reform:"bogus,pk"` // pointer field with "reform:" tag and pk label should generate error
}
