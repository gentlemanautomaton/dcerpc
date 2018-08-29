package protocol

import (
	"encoding/binary"
	"io"
)

// Tower is a protocol tower that contains binding information.
//
// Each element of a tower describes one layer of the protocol,
// known in RPC parlance as a floor. Floors are stored in
// ascending order.
//
// A particular combination of floors is called a protocol sequence.
type Tower []Floor

// WriteTo will write a binary representation of the protocol tower to w.
func (t Tower) WriteTo(w io.Writer) (n int64, err error) {
	buf := make([]byte, t.EncodedLength())
	t.Marshal(buf)
	n32, err := w.Write(buf)
	return int64(n32), err
}

// Marshal marshals the protocol tower as a binary representation stored in p.
// If len(p) is less than t.EncodedLength(), Marhsal will panic.
func (t Tower) Marshal(p []byte) {
	binary.LittleEndian.PutUint16(p[0:2], uint16(len(t)))
	offset := 2
	for i := range t {
		t[i].Marshal(p[offset:])
		offset += t[i].EncodedLength()
	}
}

// EncodedLength returns the total number of bytes required to marshal t.
func (t Tower) EncodedLength() (length int) {
	// 2-byte floor count plus the sum of all floors
	length = 2
	for i := range t {
		length += t[i].EncodedLength()
	}
	return
}
