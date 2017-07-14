package copdu

import "github.com/gentlemanautomaton/dcerpc/pdu/copdu/presentationcontext"

// BindAck represents a bind acknowledgment PDU in the connection-oriented
// protocol. It is sent from the server to the client.
type BindAck struct {
	// TODO: Decide whether we include the header here or not

	// MaxTransmitFrag is the negotiated maximum fragment size selected by the
	// server.
	MaxTransmitFrag uint16

	// MaxReceiveFrag is the negotiated maximum fragment size selected by the
	// server.
	MaxReceiveFrag uint16

	// AssocGroupID is the client-server association group that this binding is
	// associated with.
	AssocGroupID uint32

	// TODO: Add secondary address field

	// TODO: Align(4)

	// Results contains the results of the presentation context negotiation.
	Results presentationcontext.ResultList

	// TODO: Handle optional auth verifier, probably as a separate struct or something.
}
