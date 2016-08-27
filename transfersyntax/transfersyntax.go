package transfersyntax

import "github.com/gentlemanautomaton/dcerpc/uuid"

type ProtocolIdentifier []byte // TODO: Move and define. Move to core package?

// Syntax represents a DCE / RPC transfer syntax.
type Syntax interface {
	ProtocolIdentifier() ProtocolIdentifier
	UUID() uuid.UUID
}

type ndr struct {
}

var NDR = ndr{}

type ndr64 struct {
}
