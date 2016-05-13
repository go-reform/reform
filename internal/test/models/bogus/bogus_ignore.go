package bogus

//go:generate reform

// Struct without "reform:" magic comment should be ignored by ParseFile
type BogusIgnore struct {
	Bogus string `reform:"bogus"`
}
