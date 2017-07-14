package copdu

import "github.com/gentlemanautomaton/dcerpc/pdu/copdu/presentationcontext"

// Request represents a request PDU in the connection-oriented protocol.
type Request struct {
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

	// OpNum is the operation number the client would like to be invoked within
	// the RPC interface on the server.
	OpNum uint16

	// TODO: Handle optional object ID field

	// TODO: Insert stub data, 8-octet aligned

	// TODO: Handle optional auth verifier, probably as a separate struct or something.
}
