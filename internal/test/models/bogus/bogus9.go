package bogus

//go:generate reform

// Bogus9 is used for testing. reform:bogus
type Bogus9 struct {
	Bogus1 string `reform:"bogus,pk"`
	Bogus2 string `reform:"bogus"` // field with "reform:" tag with duplicate column name should generate error
}
