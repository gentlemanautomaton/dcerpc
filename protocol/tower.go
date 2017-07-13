package protocol

import "io"

// Tower is a protocol tower that contains binding information.
type Tower []Floor

// Marshal will write a binary representation of the protocol tower to the
// writer.
func (t Tower) Marshal(w io.Writer) {
	//s.Write
}
