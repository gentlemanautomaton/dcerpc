package copdu

import "github.com/gentlemanautomaton/dcerpc/formatlabel"

// Header represents the common header data shared by all connection-oriented
// protocol data units.
type Header struct {
	// TODO: Always enforce 8-octet alignment of this structure, probably through an embedded type

	// VersionMajor is the RPC protocol major version.
	VersionMajor uint8 // Should be 5
	// VersionMinor is the RPC protocol minor version.
	VersionMinor uint8
	PacketType   uint8 // 5 least significant bits
	Flags        uint8
	Format       formatlabel.Format
	FragLength   uint16
	AuthLength   uint16
	CallID       uint32
}
