package models

//go:generate reform

// types for testing
type (
	Integer int32
	String  string
	Bytes   []byte
	Uint8s  []uint8
)

//reform:extra
type Extra struct {
	ID Integer `reform:"id,pk"`
	// UUID     uuid.UUID `reform:"uuid"`
	Name *String `reform:"name"`

	Ignored1 string
	Ignored2 string `reform:""`
	Ignored3 string `reform:"-"`

	Byte    byte       `reform:"byte"`
	Uint8   uint8      `reform:"uint8"`
	ByteP   *byte      `reform:"bytep"`
	Uint8P  *uint8     `reform:"uint8p"`
	Bytes   []byte     `reform:"bytes"`
	Uint8s  []uint8    `reform:"uint8s"`
	BytesA  [512]byte  `reform:"bytesa"`
	Uint8sA [512]uint8 `reform:"uint8sa"`
	BytesT  Bytes      `reform:"bytest"`
	Uint8sT Uint8s     `reform:"uint8st"`
}

//reform:not_exported
type notExported struct {
	ID string `reform:"id,pk"`
}
