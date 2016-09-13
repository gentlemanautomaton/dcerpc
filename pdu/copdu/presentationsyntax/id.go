package presentationsyntax

import "github.com/nu7hatch/gouuid"

// ID contains the interface UUID and version of a presentation syntax.
type ID struct {
	Interface uuid.UUID
	Version   uint32
}
