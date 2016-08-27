package types

import "github.com/gentlemanautomaton/dcerpc/uuid"

// Interface represents an Interface as expressed within the DCE / RPC interface
// definition language.
type Interface struct {
	uuid uuid.UUID
}
