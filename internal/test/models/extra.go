package models

//go:generate reform

// types for testing
type (
	Integer int32
	String  string
	Bytes   []byte
)

//reform:extra
type Extra struct {
	ID     Integer   `reform:"id,pk"`
	Name   *String   `reform:"name"`
	Bytes  []byte    `reform:"bytes"`
	Bytes2 Bytes     `reform:"bytes2"`
	Byte   *byte     `reform:"byte"`
	Array  [512]byte `reform:"array"`
}

//reform:not_exported
type notExported struct {
	ID string `reform:"id,pk"`
}
