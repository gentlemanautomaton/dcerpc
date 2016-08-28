package formatlabel

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
