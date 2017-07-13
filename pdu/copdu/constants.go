package copdu

// Connection-oriented Packet Flags
const (
	// FirstFrag indicates that the PDU is the first fragment in a series.
	FirstFrag = 0x01
	// LastFrag indicates that the PDU is the last frament in a series.
	LastFrag = 0x02
	// PendingCancel indicates that a cancellation was pending when the PDU was
	// sent.
	PendingCancel = 0x04
	// ConcurrentMultiplexing indicates that the underlying connection supports
	// concurrent multiplexing.
	ConcurrentMultiplexing = 0x10
	// DidNotExecute is only used in fault packets, and indicates that the remote
	// procedure did not execute.
	DidNotExecute = 0x20
	// Maybe indicates that "maybe" semantics have been requested.
	Maybe = 0x40
	// ObjectUUID indicates that a non-nil UUID is present in the object field.
	// If this flag has not been set the object field has not been included.
	ObjectUUID = 0x80
)
