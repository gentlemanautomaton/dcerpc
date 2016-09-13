package presentationcontext

// Reason represents the reason that a presentation context was rejected
// by a provider.
type Reason uint16

const (
	// ReasonNotSpecified indicates that a presentation context negotiation
	// failed for an unspecified reason.
	ReasonNotSpecified Reason = iota

	// AbstractSyntaxNotSupported indicates that a presentation context
	// negotiation failed because the abstract transfer syntax was not supported.
	AbstractSyntaxNotSupported

	// ProposedTransferSyntaxesNotSupported indicates that a presentation context
	// negotiation failed because none of the proposed transfer syntaxes are
	// supported.
	ProposedTransferSyntaxesNotSupported

	// LocalLimitExceeded indicates that a presentation context negotiation
	// failed because a local limit was exceeded by the other party.
	LocalLimitExceeded
)
