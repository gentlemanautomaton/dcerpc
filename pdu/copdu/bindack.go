package copdu

import "github.com/gentlemanautomaton/dcerpc/pdu/copdu/presentationcontext"

// BindAck represents a bind acknowledgment PDU in the connection-oriented
// protocol. It is sent from the server to the client.
type BindAck struct {
	// TODO: Decide whether we include the header here or not

	// MaxTransmitFrag is the maximum fragment size the sender can transmit.
	MaxTransmitFrag uint16

	// MaxReceiveFrag is the maximum fragment size the sender can receive.
	MaxReceiveFrag uint16

	// AssocGroupID is the client-server association group that this packet is
	// associated with. 0 indicates a request for the creation of a new group.
	AssocGroupID uint32

	// TODO: Add secondary address field

	// TODO: Align(4)

	Results presentationcontext.ResultList

	// TODO: Handle optional auth verifier, probably as a separate struct or something.
}
