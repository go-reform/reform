package bogus

//go:generate reform

// Bogus11 is used for testing. reform:bogus
type Bogus11 struct {
	Bogus []string `reform:"bogus,pk"` // slice field with "reform:" tag and pk label should generate error
}
