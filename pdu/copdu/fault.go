package copdu

import "github.com/gentlemanautomaton/dcerpc/pdu/copdu/presentationcontext"

// Fault represents a fault PDU in the connection-oriented protocol.
type Fault struct {
	// TODO: Decide whether we include the header here or not

	AllocHint             uint32
	PresentationContextID presentationcontext.ID
	CancelCount           uint8
	_                     uint8 // Reserved
	Status                uint32
	_                     [4]uint8 // Reserved / 8-octet alignment

	// TODO: Add stub data, 8-octet aligned

	// TODO: Handle optional auth verifier, probably as a separate struct or something.
}
