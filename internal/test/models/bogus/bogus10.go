package bogus

//go:generate reform

// Bogus10 is used for testing. reform:bogus
type Bogus10 struct {
	Bogus1 string `reform:"bogus1,pk"`
	Bogus2 string `reform:"bogus2,pk"` // field with "reform:" tag with duplicate pk label should generate error
}
