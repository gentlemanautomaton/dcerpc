package ndr

import (
	"errors"
	"io"

	"github.com/gentlemanautomaton/dcerpc/formatlabel"
)

const readerBufLen = 16

// A Reader is capable of reading all NDR primitive types.
type Reader interface {
	Offset() uint64
	Skip(count int) (err error)
	Align(modulo int) (err error)
	Read(p []byte) (n int, err error)
	ReadFull(p []byte) (err error)
	ReadByte() (v byte, err error)
	ReadBool() (v bool, err error)
	ReadASCII() (v string, err error) // FIXME: Decide whether these should write individual runes instead
	ReadEBCDIC() (v string, err error)
	ReadUnicode() (v string, err error)
	ReadInt8() (v int8, err error)
	ReadUint8() (v uint8, err error)

	// Big-endian integer representations

	ReadInt16BE() (v int16, err error)
	ReadInt32BE() (v int32, err error)
	ReadInt64BE() (v int64, err error)
	ReadUint16BE() (v uint16, err error)
	ReadUint32BE() (v uint32, err error)
	ReadUint64BE() (v uint64, err error)

	// Little-endian integer representations

	ReadInt16LE() (v int16, err error)
	ReadInt32LE() (v int32, err error)
	ReadInt64LE() (v int64, err error)
	ReadUint16LE() (v uint16, err error)
	ReadUint32LE() (v uint32, err error)
	ReadUint64LE() (v uint64, err error)

	// IEEE floating point representations
	ReadFloat32BEIEEE() (v float32, err error)
	ReadFloat64BEIEEE() (v float64, err error)
	ReadFloat32LEIEEE() (v float32, err error)
	ReadFloat64LEIEEE() (v float64, err error)

	// TODO: Add VAX, Cray and IBM floating point representations

	// Format-dependent string representations
	ReadString() (v string, err error) // FIXME: Change to ReadCharacter instead?

	// Format-dependent integer representations

	ReadInt16() (v int16, err error)
	ReadInt32() (v int32, err error)
	ReadInt64() (v int64, err error)

	ReadUint16() (v uint16, err error)
	ReadUint32() (v uint32, err error)
	ReadUint64() (v uint64, err error)

	// TODO: Add Enums?

	// Format-dependent floating point representations
	ReadFloat32() (v float32, err error)
	ReadFloat64() (v float64, err error)

	// TODO: Add referent recording funtion
	// ReadPointer()
}

// NewReader returns a new Reader that will read from the given underlying
// io.Reader using the given format label.
func NewReader(r io.Reader, format formatlabel.Format) Reader {
	switch format {
	case formatlabel.BEAIEEE:
		return &readerBEAIEEE{reader{
			Reader: r,
			refs:   make(map[uintptr]uint64),
		}}
	}
	return nil
}

// reader provides all of NDR primitive reading functions except for those
// that are format-dependent. Data is written to an underlying io.Reader.
type reader struct {
	io.Reader
	index uint64 // Current offset from the start of the octet stream
	buf   [readerBufLen]byte
	// NOTE: If Go ever implements a compacting GC it will be important that we
	//       pin any of the pointers we are tracking here.
	nextRefID uint64             // ID of the next referent; must be > 0
	refs      map[uintptr]uint64 // Referent map translates pointers to referent IDs of already-encoded referents (structures)
}

func (r *reader) Offset() uint64 {
	return r.index
}

func (r *reader) Skip(count int) (err error) {
	// Small reads get the fast path
	if count <= readerBufLen {
		return r.ReadFull(r.buf[0:count])
	}

	// Attempt to use an appropriately sized buffer in the hope that it fits on
	// the stack.
	switch {
	case count <= 64:
		var buf [64]byte
		return r.ReadFull(buf[0:count])
	case count <= 128:
		var buf [128]byte
		return r.ReadFull(buf[0:count])
	case count <= 256:
		var buf [256]byte
		return r.ReadFull(buf[0:count])
	case count <= 512:
		var buf [512]byte
		return r.ReadFull(buf[0:count])
	}

	// Large reads get the slow path
	const bufLen = 4096

	var buf [bufLen]byte
	for count > 0 && err == nil {
		chunk := count
		if chunk > bufLen {
			chunk = bufLen
		}
		err = r.ReadFull(buf[0:chunk])
		count -= chunk
	}

	return
}

// Align will read until the current index of the buffer matches the desired
// alignment (until index mod modulo == zero).
func (r *reader) Align(modulo int) (err error) {
	m := int(r.index % uint64(modulo))
	if m == 0 {
		return
	}

	return r.Skip(modulo - m)
}

// Alloc will ensure that the given number of bytes have been preallocated
// within the internal buffer.
/*
func (r *reader) Alloc(bytes uint64) {
	//c := cap(s.b.data)
	//if c
}
*/

func (r *reader) Read(buf []byte) (n int, err error) {
	n, err = r.Reader.Read(buf)
	r.index += uint64(n)
	return
}

// ReadFull reads exactly len(buf) bytes from the reader into buf.
func (r *reader) ReadFull(buf []byte) (err error) {
	n, needed := 0, len(buf)
	for n < needed && err == nil {
		var nn int
		nn, err = r.Reader.Read(buf)
		n += nn
	}
	r.index += uint64(n)
	if n == needed {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

func (r *reader) ReadByte() (v byte, err error) {
	var n int
	n, err = r.Reader.Read(r.buf[0:1])
	if n == 1 {
		r.index++
		v = r.buf[0]
	}
	return
}

func (r *reader) ReadBool() (v bool, err error) {
	b, err := r.ReadByte()
	if err != nil {
		return
	}
	if b == 0 {
		return false, nil
	}
	return true, nil
}

// TODO: Determine whether this should read individual runes instead?
func (r *reader) ReadASCII() (v string, err error) {
	// FIXME: Implement this
	return
}

// TODO: Determine whether this should read individual runes instead?
func (r *reader) ReadEBCDIC() (v string, err error) {
	// FIXME: Implement this
	return
}

// TODO: Determine whether this should read individual runes instead?
func (r *reader) ReadUnicode() (v string, err error) {
	// FIXME: Implement this
	return
}

func (r *reader) ReadInt8() (v int8, err error) {
	b, err := r.ReadByte()
	if err != nil {
		return
	}
	return int8(b), nil
}

func (r *reader) ReadUint8() (v uint8, err error) {
	b, err := r.ReadByte()
	if err != nil {
		return
	}
	return uint8(b), nil
}

func (r *reader) ReadInt16BE() (v int16, err error) {
	u, err := r.ReadUint16BE()
	return int16(u), err
}

func (r *reader) ReadInt32BE() (v int32, err error) {
	u, err := r.ReadUint32BE()
	return int32(u), err
}

func (r *reader) ReadInt64BE() (v int64, err error) {
	u, err := r.ReadUint64BE()
	return int64(u), err
}

func (r *reader) ReadUint16BE() (v uint16, err error) {
	r.Align(2)
	if err = r.ReadFull(r.buf[0:2]); err != nil {
		return
	}
	v = (uint16(r.buf[0]) << 8) | uint16(r.buf[1])
	return
}

func (r *reader) ReadUint32BE() (v uint32, err error) {
	r.Align(4)
	if err = r.ReadFull(r.buf[0:4]); err != nil {
		return
	}
	v = (uint32(r.buf[0]) << 24) |
		(uint32(r.buf[1]) << 16) |
		(uint32(r.buf[2]) << 8) |
		uint32(r.buf[3])
	return
}

func (r *reader) ReadUint64BE() (v uint64, err error) {
	r.Align(8)
	if err = r.ReadFull(r.buf[0:4]); err != nil {
		return
	}
	v = (uint64(r.buf[0]) << 56) |
		(uint64(r.buf[1]) << 48) |
		(uint64(r.buf[2]) << 40) |
		(uint64(r.buf[3]) << 32) |
		(uint64(r.buf[4]) << 24) |
		(uint64(r.buf[5]) << 16) |
		(uint64(r.buf[6]) << 8) |
		uint64(r.buf[7])
	return
}

func (r *reader) ReadInt16LE() (v int16, err error) {
	u, err := r.ReadUint16LE()
	return int16(u), err
}

func (r *reader) ReadInt32LE() (v int32, err error) {
	u, err := r.ReadUint32LE()
	return int32(u), err
}

func (r *reader) ReadInt64LE() (v int64, err error) {
	u, err := r.ReadUint64LE()
	return int64(u), err
}

func (r *reader) ReadUint16LE() (v uint16, err error) {
	r.Align(2)
	if err = r.ReadFull(r.buf[0:2]); err != nil {
		return
	}
	v = uint16(r.buf[0]) | (uint16(r.buf[1]) << 8)
	return
}

func (r *reader) ReadUint32LE() (v uint32, err error) {
	r.Align(4)
	if err = r.ReadFull(r.buf[0:4]); err != nil {
		return
	}
	v = uint32(r.buf[0]) |
		(uint32(r.buf[1]) << 8) |
		(uint32(r.buf[2]) << 16) |
		(uint32(r.buf[3]) << 24)
	return
}

func (r *reader) ReadUint64LE() (v uint64, err error) {
	r.Align(8)
	if err = r.ReadFull(r.buf[0:4]); err != nil {
		return
	}
	v = uint64(r.buf[0]) |
		(uint64(r.buf[1]) << 8) |
		(uint64(r.buf[2]) << 16) |
		(uint64(r.buf[3]) << 24) |
		(uint64(r.buf[4]) << 32) |
		(uint64(r.buf[5]) << 40) |
		(uint64(r.buf[6]) << 48) |
		(uint64(r.buf[7]) << 56)
	return
}

func (r *reader) ReadFloat32BEIEEE() (v float32, err error) {
	// TODO: Write this
	return 0, errors.New("Not implemented")
}

func (r *reader) ReadFloat64BEIEEE() (v float64, err error) {
	// TODO: Write this
	return 0, errors.New("Not implemented")
}

func (r *reader) ReadFloat32LEIEEE() (v float32, err error) {
	// TODO: Write this
	return 0, errors.New("Not implemented")
}

func (r *reader) ReadFloat64LEIEEE() (v float64, err error) {
	// TODO: Write this
	return 0, errors.New("Not implemented")
}

var _ = Reader((*readerBEAIEEE)(nil)) // Compile-time check for interface compliance

// readerBEAIEEE writes NDR-encoded primitive data types with big-endian
// integer representation, ASCII character representation and IEEE floating
// point representation.
type readerBEAIEEE struct {
	reader
}

func (r *readerBEAIEEE) ReadString() (v string, err error) {
	return r.ReadASCII()
}

func (r *readerBEAIEEE) ReadInt16() (v int16, err error) {
	return r.ReadInt16BE()
}

func (r *readerBEAIEEE) ReadInt32() (v int32, err error) {
	return r.ReadInt32BE()
}

func (r *readerBEAIEEE) ReadInt64() (v int64, err error) {
	return r.ReadInt64BE()
}

func (r *readerBEAIEEE) ReadUint16() (v uint16, err error) {
	return r.ReadUint16BE()
}

func (r *readerBEAIEEE) ReadUint32() (v uint32, err error) {
	return r.ReadUint32BE()
}

func (r *readerBEAIEEE) ReadUint64() (v uint64, err error) {
	return r.ReadUint64BE()
}

func (r *readerBEAIEEE) ReadFloat32() (v float32, err error) {
	return r.ReadFloat32BEIEEE()
}

func (r *readerBEAIEEE) ReadFloat64() (v float64, err error) {
	return r.ReadFloat64BEIEEE()
}
