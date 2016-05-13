package bogus

//go:generate reform

// reform:bogus
type Bogus6 struct {
	// struct without fields with "reform:" tag should generate error
	Bogus string
}
