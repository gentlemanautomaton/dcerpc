package presentationcontext

import "github.com/gentlemanautomaton/dcerpc/pdu/copdu/presentationsyntax"

// Element represents a presentation context element.
//
// TODO: Consider renaming this to RequestElement.
type Element struct {
	ID                  ID
	NumTransferSyntaxes uint8
	_                   uint8 // Reserved
	AbstractSyntax      presentationsyntax.ID
	TransferSyntaxes    []presentationsyntax.ID `idl:"size_is(NumTransferSyntaxes)"`
}
