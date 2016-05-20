package bogus

//go:generate reform

// Bogus4 is used for testing. reform:bogus
type Bogus4 struct {
	Bogus string `reform:",pk"` // field with "reform:" tag without column name should generate error
}
