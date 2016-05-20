package bogus

//go:generate reform

// Bogus6 is used for testing. reform:bogus
type Bogus6 struct {
	// struct without fields with "reform:" tag should generate error
	Bogus string
}
