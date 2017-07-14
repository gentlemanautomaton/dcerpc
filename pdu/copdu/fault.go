package copdu

import "github.com/gentlemanautomaton/dcerpc/pdu/copdu/presentationcontext"

// Fault represents a fault PDU in the connection-oriented protocol.
type Fault struct {
	// TODO: Decide whether we include the header here or not

	// AllocHint is an optional suggested buffer size provided by the sender of
	// a fragmented PDU series. When used, it indiciates the amount of memory
	// required to hold the entire series of fragmented requests in a contiguous
	// block. When no hint is provided AllocHint will be zero.
	AllocHint uint32

	// PresContextID is the presentation context identifier of the selected
	// syntax. The identifier is chosen by the client when assembling its
	// supported presentation context list.
	PresContextID presentationcontext.ID

	// CancelCount is the number of cancellations that have been received by the
	// server.
	CancelCount uint8

	_ uint8 // Reserved

	// Status describes the nature of the fault condition. If Status is nonzero
	// it represents a standard RPC runtime error code; the PDU must not contain
	// stub data in this case. If Status is zero it is an application error that
	// will be described in the stub data; the manner of its description is
	// application-specific.
	Status uint32 // Fault Code

	_ [4]uint8 // Reserved / 8-octet alignment

	// TODO: Add stub data, 8-octet aligned

	// TODO: Handle optional auth verifier, probably as a separate struct or something.
}
