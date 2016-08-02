package ndr

import "io"

const zeroPaddingLen = 512
const writerBufLen = 16

var zeroPadding [zeroPaddingLen]byte

// A Writer is capable of writing all NDR primitive types
type Writer interface {
	Offset() uint64
	Align(modulo int)
	Write(p []byte) (int, error)
	WriteByte(c byte) error
	WriteBool(v bool)
	WriteASCII(v string) // FIXME: Decide whether these should write individual runes instead
	WriteEBCDIC(v string)
	WriteUnicode(v string)
	WriteInt8(v int8)
	WriteUint8(v uint8)

	// Big-endian integer representations

	WriteInt16BE(v int16)
	WriteInt32BE(v int32)
	WriteInt64BE(v int64)
	WriteUint16BE(v uint16)
	WriteUint32BE(v uint32)
	WriteUint64BE(v uint64)

	// Little-endian integer representations

	WriteInt16LE(v int16)
	WriteInt32LE(v int32)
	WriteInt64LE(v int64)
	WriteUint16LE(v uint16)
	WriteUint32LE(v uint32)
	WriteUint64LE(v uint64)

	// IEEE floating point representations
	WriteFloat32BEIEEE(v float32)
	WriteFloat64BEIEEE(v float64)
	WriteFloat32LEIEEE(v float32)
	WriteFloat64LEIEEE(v float64)

	// TODO: Add VAX, Cray and IBM floating point representations

	// Format-dependent string representations
	WriteString(v string) // FIXME: Change to WriteCharacter instead?

	// Format-dependent integer representations

	WriteInt16(v int16)
	WriteInt32(v int32)
	WriteInt64(v int64)

	WriteUint16(v uint16)
	WriteUint32(v uint32)
	WriteUint64(v uint64)

	// TODO: Add Enums?

	// Format-dependent floating point representations
	WriteFloat32(v float32)
	WriteFloat64(v float64)

	// TODO: Add referent recording funtion
	// WritePointer()
}

// NewWriter returns a new Writer that will write to the given underlying
// io.Writer using the given format label.
func NewWriter(w io.Writer, format FormatLabel) Writer {
	switch {
	case format.CharacterRepresentation() == ASCII && format.IntegerRepresentation() == BigEndian && format.FloatRepresentation() == IEEE:
		return &writerBEAIEEE{writer{
			Writer: w,
			refs:   make(map[uintptr]uint64),
		}}
	}
	return nil
}

// writer provides all of NDR primitive writing functions except for those
// that are format-dependent. Data is written to an underlying io.Writer.
type writer struct {
	io.Writer
	index uint64 // Current offset from the start of the octet stream
	buf   [writerBufLen]byte
	// NOTE: If Go ever implements a compacting GC it will be important that we
	//       pin any of the pointers we are tracking here.
	nextRefID uint64             // ID of the next referent; must be > 0
	refs      map[uintptr]uint64 // Referent map translates pointers to referent IDs of already-encoded referents (structures)
}

func (w *writer) Offset() uint64 {
	return w.index
}

// Align will insert zero padding until the current index of the buffer
// matches the desired alignment (until index mod modulo == zero).
func (w *writer) Align(modulo int) {
	m := int(w.index % uint64(modulo))
	if m == 0 {
		return
	}

	remaining := modulo - m // number of zero bytes to write
	if remaining <= zeroPaddingLen {
		w.Write(zeroPadding[0:remaining])
		return
	}

	for remaining > 0 {
		chunk := remaining
		if chunk > zeroPaddingLen {
			chunk = zeroPaddingLen
		}
		w.Write(zeroPadding[0:chunk])
		remaining -= chunk
	}
}

func (w *writer) WriteByte(c byte) error {
	w.buf[0] = c
	_, err := w.Write(w.buf[0:1])
	if err != nil {
		return err
	}
	w.index++
	return nil
}

// Alloc will ensure that the given number of bytes have been preallocated
// within the internal buffer.
/*
func (w *writer) Alloc(bytes uint64) {
	//c := cap(s.b.data)
	//if c
}
*/

func (w *writer) WriteBool(v bool) {
	if v {
		w.WriteByte(1)
	} else {
		w.WriteByte(0)
	}
}

// TODO: Determine whether this should write individual runes instead?
func (w *writer) WriteASCII(v string) {
	// FIXME: Find a better way of ASCII-encoding UTF-8 characters, or return an error
	w.Write([]byte(v))
}

// TODO: Determine whether this should write individual runes instead?
func (w *writer) WriteEBCDIC(v string) {
	// FIXME: Implement this
}

// TODO: Determine whether this should write individual runes instead?
func (w *writer) WriteUnicode(v string) {
	// FIXME: Implement this
}

func (w *writer) WriteInt8(v int8) {
	w.WriteByte(byte(v))
}

func (w *writer) WriteUint8(v uint8) {
	w.WriteByte(byte(v))
}

func (w *writer) WriteInt16BE(v int16) {
	w.WriteUint16BE(uint16(v))
}

func (w *writer) WriteInt32BE(v int32) {
	w.WriteUint32BE(uint32(v))
}

func (w *writer) WriteInt64BE(v int64) {
	w.WriteUint64BE(uint64(v))
}

func (w *writer) WriteUint16BE(v uint16) {
	w.Align(2)
	w.buf[0] = byte(v >> 8)
	w.buf[1] = byte(v)
	w.Write(w.buf[0:2])
}

func (w *writer) WriteUint32BE(v uint32) {
	w.Align(4)
	w.buf[0] = byte(v >> 24)
	w.buf[1] = byte(v >> 16)
	w.buf[2] = byte(v >> 8)
	w.buf[3] = byte(v)
	w.Write(w.buf[0:4])
}

func (w *writer) WriteUint64BE(v uint64) {
	w.Align(8)
	w.buf[0] = byte(v >> 56)
	w.buf[1] = byte(v >> 48)
	w.buf[2] = byte(v >> 40)
	w.buf[3] = byte(v >> 32)
	w.buf[4] = byte(v >> 24)
	w.buf[5] = byte(v >> 16)
	w.buf[6] = byte(v >> 8)
	w.buf[7] = byte(v)
	w.Write(w.buf[0:8])
}

func (w *writer) WriteInt16LE(v int16) {
	w.WriteUint16LE(uint16(v))
}

func (w *writer) WriteInt32LE(v int32) {
	w.WriteUint32LE(uint32(v))
}

func (w *writer) WriteInt64LE(v int64) {
	w.WriteUint64LE(uint64(v))
}

func (w *writer) WriteUint16LE(v uint16) {
	w.Align(2)
	w.buf[0] = byte(v)
	w.buf[1] = byte(v >> 8)
	w.Write(w.buf[0:2])
}

func (w *writer) WriteUint32LE(v uint32) {
	w.Align(4)
	w.buf[0] = byte(v)
	w.buf[1] = byte(v >> 8)
	w.buf[2] = byte(v >> 16)
	w.buf[3] = byte(v >> 24)
	w.Write(w.buf[0:4])
}

func (w *writer) WriteUint64LE(v uint64) {
	w.Align(8)
	w.buf[0] = byte(v)
	w.buf[1] = byte(v >> 8)
	w.buf[2] = byte(v >> 16)
	w.buf[3] = byte(v >> 24)
	w.buf[4] = byte(v >> 32)
	w.buf[5] = byte(v >> 40)
	w.buf[6] = byte(v >> 48)
	w.buf[7] = byte(v >> 56)
	w.Write(w.buf[0:8])
}

func (w *writer) WriteFloat32BEIEEE(v float32) {
	// TODO: Write this
}

func (w *writer) WriteFloat64BEIEEE(v float64) {
	// TODO: Write this
}

func (w *writer) WriteFloat32LEIEEE(v float32) {
	// TODO: Write this
}

func (w *writer) WriteFloat64LEIEEE(v float64) {
	// TODO: Write this
}

var _ = Writer((*writerBEAIEEE)(nil)) // Compile-time check for interface compliance

// writerBEAIEEE writes NDR-encoded primitive data types with big-endian
// integer representation, ASCII character representation and IEEE floating
// point representation.
type writerBEAIEEE struct {
	writer
}

func (w *writerBEAIEEE) WriteString(v string) {
	w.WriteASCII(v)
}

func (w *writerBEAIEEE) WriteInt16(v int16) {
	w.WriteInt16BE(v)
}

func (w *writerBEAIEEE) WriteInt32(v int32) {
	w.WriteInt32BE(v)
}

func (w *writerBEAIEEE) WriteInt64(v int64) {
	w.WriteInt64BE(v)
}

func (w *writerBEAIEEE) WriteUint16(v uint16) {
	w.WriteUint16BE(v)
}

func (w *writerBEAIEEE) WriteUint32(v uint32) {
	w.WriteUint32BE(v)
}

func (w *writerBEAIEEE) WriteUint64(v uint64) {
	w.WriteUint64BE(v)
}

func (w *writerBEAIEEE) WriteFloat32(v float32) {
	w.WriteFloat32BEIEEE(v)
}

func (w *writerBEAIEEE) WriteFloat64(v float64) {
	w.WriteFloat64BEIEEE(v)
}
