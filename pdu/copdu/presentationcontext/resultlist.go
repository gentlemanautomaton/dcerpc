package presentationcontext

// ResultList represents the list of results of a presentation context
// negotiation.
type ResultList struct {
	NumResults uint8
	_          uint8           // Reserved
	_          uint16          // Reserved
	Results    []ResultElement `idl:"size_is(NumResults)"`
}
