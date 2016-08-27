package formatlabel

//type IntegerRepresentation byte

//type CharacterRepresentation byte

const (
	// ASCII represents ASCII-encoded text characters.
	ASCII = 0
	// EBCDIC represents EBCDIC-encoded text characters.
	EBCDIC = 1
)

const (
	// BigEndian represents big-endian encoding of mutli-byte characters.
	BigEndian = 0
	// LittleEndian represents little-endian encoding of mutli-byte characters.
	LittleEndian = 1
)

const (
	// IEEE represents the IEEE floating point representation.
	IEEE = 0
	// VAX represents the VAX floating point representation.
	VAX = 1
	// Cray represents the Cray floating point representation.
	Cray = 2
	// IBM represents the IBM floating point representation.
	IBM = 3
)

var (
	// BEAIEEE is the format label for big-endian, ASCII and IEEE encoding.
	BEAIEEE = New(BigEndian, ASCII, IEEE)
	// LEAIEEE is the format label for little-endian, ASCII and IEEE encoding.
	LEAIEEE = New(LittleEndian, ASCII, IEEE)
)

// Format represents the NDR encoding format
type Format [4]byte

// New returns a format label for the given integer, character and
// floating point representations.
func New(intRep, charRep, floatRep byte) Format {
	return [4]byte{((intRep & 0x0f) << 4) | (charRep & 0x0f), floatRep, 0, 0}
}

// IntegerRepresentation returns the identifier of the integer representation
// of the format label.
func (f *Format) IntegerRepresentation() byte {
	return byte(f[0] >> 4)
}

// CharacterRepresentation returns the identifier of the character
// representation of the format label.
func (f *Format) CharacterRepresentation() byte {
	return f[0] & 0x00ff
}

// FloatRepresentation returns the identifier of the floating-point
// representation of the format label.
func (f *Format) FloatRepresentation() byte {
	return f[0] & 0x00ff
}
