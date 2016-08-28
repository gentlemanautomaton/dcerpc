package formatlabel

// Format represents the NDR encoding format
type Format [4]byte

// New returns a format label for the given integer, character and
// floating point representations.
func New(intRep, charRep, floatRep byte) Format {
	return [4]byte{((intRep & 0x0f) << 4) | (charRep & 0x0f), floatRep, 0, 0}
}

// IntRep returns the integer representation of the format label.
func (f *Format) IntRep() byte {
	return f[0] >> 4
}

// CharRep returns the character representation of the format label.
func (f *Format) CharRep() byte {
	return f[0] & 0x0f
}

// FloatRep returns the the floating-point representation of the format label.
func (f *Format) FloatRep() byte {
	return f[1]
}
