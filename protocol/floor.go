package protocol

import (
	"encoding/binary"
	"io"
)

// Floor represents a floor in the protocol tower. It describes one
// layer of the protocol stack.
type Floor struct {
	ProtocolIdentifier []byte
	AddressData        []byte
}

// WriteTo will write a binary representation of the protocol floor to w.
func (f Floor) WriteTo(w io.Writer) (n int64, err error) {
	buf := make([]byte, f.EncodedLength())
	f.Marshal(buf)
	n32, err := w.Write(buf)
	return int64(n32), err
}

// Marshal marshals the protocol floor as a binary representation stored in p.
// If len(p) is less than f.EncodedLength(), Marhsal will panic.
func (f Floor) Marshal(p []byte) {
	copyWithLength(f.ProtocolIdentifier, p) // Left-hand side (key)
	p = p[2+len(f.ProtocolIdentifier):]
	copyWithLength(f.AddressData, p) // Right-hand side (value)
}

// EncodedLength returns the total number of bytes required to encode f.
func (f Floor) EncodedLength() int {
	return 4 + len(f.ProtocolIdentifier) + len(f.AddressData)
}

func copyWithLength(from []byte, to []byte) {
	binary.LittleEndian.PutUint16(to[0:2], uint16(len(from)))
	copy(to[2:], from)
}
