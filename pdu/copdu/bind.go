package copdu

import "github.com/gentlemanautomaton/dcerpc/pdu/copdu/presentationcontext"

// Bind represents a bind PDU in the connection-oriented protocol. It is sent
// from the client to the server.
type Bind struct {
	// TODO: Decide whether we include the header here or not

	// MaxTransmitFrag is the maximum fragment size the client would like to
	// transmit.
	MaxTransmitFrag uint16

	// MaxReceiveFrag is the maximum fragment size the client would like to
	// receive.
	MaxReceiveFrag uint16

	// AssocGroupID is the client-server association group that this binding is
	// associated with. 0 indicates a request for the creation of a new group.
	AssocGroupID uint32

	// Elements is a variable-length ordered list of supported presentation
	// syntaxes that the client is offering for negotiation.
	Elements presentationcontext.List

	// TODO: Handle optional auth verifier, probably as a separate struct or something.
}
