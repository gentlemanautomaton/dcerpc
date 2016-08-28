package copdu

// Connection-oriented Packet Flags
const (
	// FirstFrag

	// LastFrag indicates that the PDU is the last frament in a series.
	LastFrag = 0x02
	// Frag indicates that the PDU is one fragment within a series.
	Frag = 0x04
	// No
)
