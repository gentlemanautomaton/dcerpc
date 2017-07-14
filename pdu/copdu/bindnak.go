package copdu

import "github.com/gentlemanautomaton/dcerpc/pdu/copdu/presentationcontext"

// BindNak represents a bind rejection PDU in the connection-oriented
// protocol. It is sent from the server to the client.
type BindNak struct {
	// TODO: Decide whether we include the header here or not

	// RejectReason indicates why the binding was rejected.
	RejectReason presentationcontext.Reason

	// TODO: Add array of protocol versions supported
}
