package presentationcontext

// List represents a presentation context list.
//
// TODO: Consider renaming this to RequestList.
type List struct {
	NumElements uint8
	_           uint8     // Reserved
	_           uint16    // Reserved
	Elements    []Element `idl:"size_is(NumElements)"`
}
