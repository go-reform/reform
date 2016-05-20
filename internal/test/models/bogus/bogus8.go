package bogus

//go:generate reform

// Bogus8 is used for testing. reform:bogus
type Bogus8 struct {
	Bogus *string `reform:"bogus,omitempty"` // pointer field with "reform:" tag and omitempty label should generate error
}
