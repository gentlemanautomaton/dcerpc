package pdu

// Packet types
const (
	// TypeRequest indicates a request packet in both connection-oriented and
	// connectionless protocols.
	TypeRequest = iota

	// TypePing indicates a ping packet in the connectionless protocol.
	TypePing

	// TypeResponse indicates a response packet in both connection-oriented and
	// connectionless protocols.
	TypeResponse

	// TypeFault indicates a fault packet in both connection-oriented and
	// connectionless protocols.
	TypeFault

	// TypeWorking indicates a working packet in the connectionless protocol.
	TypeWorking

	// TypeNoCall indicates a nocall packet in the connectionless protocol.
	TypeNoCall

	// TypeReject indicates a reject packet in the connectionless protocol.
	TypeReject

	// TypeAck indicates an ack packet in the connectionless protocol.
	TypeAck

	// TypeCancelCL indicates a cl_cancel packet in the connectionless protocol.
	TypeCancelCL

	// TypeFAck indicates an fack packet in the connectionless protocol.
	TypeFAck

	// TypeCancelAck indicates a cancel_ack packet in the connectionless protocol.
	TypeCancelAck

	// TypeBind indicates a bind packet in the connection-oriented protocol.
	TypeBind

	// TypeBindAck indicates a bindack packet in the connection-oriented
	// protocol.
	TypeBindAck

	// TypeBindNak indicates a bindnack packet in the connection-oriented
	// protocol.
	TypeBindNak

	// TypeAlterContext indicates an alter_context packet in the
	// connection-oriented protocol.
	TypeAlterContext

	// TypeAlterContextResp indicates an alter_context_resp packet in the
	// connection-oriented protocol.
	TypeAlterContextResp

	// TypeShutdown indicates a shutdown packet in the connection-oriented
	// protocol.
	TypeShutdown

	// TypeCancelCO indicates a cancel packet in the connection-oriented protocol.
	TypeCancelCO
	// TypeOrphaned indicates an orphaned packet in the connection-oriented
	// protocol.
	TypeOrphaned
)
