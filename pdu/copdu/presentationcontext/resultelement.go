package presentationcontext

import "github.com/gentlemanautomaton/dcerpc/pdu/copdu/presentationsyntax"

// ResultElement represents the result of a presentation context negotiation.
type ResultElement struct {
	// Result indicates the server's response to the proposed presentation
	// context.
	Result Result
	// Reason indicates the reason the presentation context was rejected in the
	// case of provider rejection.
	Reason Reason
	// TransferSyntax holds the selected transfer syntax in the case of
	// acceptance, otherwise it will be zero.
	TransferSyntax presentationsyntax.ID
}
