package ndr

//type IntegerRepresentation byte

//type CharacterRepresentation byte

const (
	ASCII  = 0
	EBCDIC = 1
	// Unicode?
)

const (
	BigEndian    = 0
	LittleEndian = 1
)

const (
	IEEE = 0
	VAX  = 1
	Cray = 2
	IBM  = 3
)

// FormatLabel represents the NDR encoding selections
type FormatLabel [4]byte

func (f *FormatLabel) IntegerRepresentation() byte {
	return byte(f[0] >> 4)
}

func (f *FormatLabel) CharacterRepresentation() byte {
	return f[0] & 0x00ff
}

func (f *FormatLabel) FloatRepresentation() byte {
	return f[0] & 0x00ff
}
