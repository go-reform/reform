package bogus

//go:generate reform

// reform:bogus
type Bogus4 struct {
	Bogus string `reform:",pk"` // field with "reform:" tag without column name should generate error
}
