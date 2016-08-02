package dcerpc

import "io"

type Tower []Floor

type Floor struct {
	ProtocolIdentifier []byte
	AddressData        []byte
}

func (t Tower) Marshal(w io.Writer) {
	//s.Write
}
