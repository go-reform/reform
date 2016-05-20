package bogus

//go:generate reform

// BogusIgnore is used for testing.
// Struct without "reform:" magic comment should be ignored by ParseFile.
type BogusIgnore struct {
	Bogus string `reform:"bogus"`
}
