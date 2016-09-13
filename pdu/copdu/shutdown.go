package copdu

// Shutdown represents a shutdown PDU in the connection-oriented protocol. It is
// sent from the server to the client to request that the client end the
// connection and release any allocated resources.
type Shutdown struct {
	// TODO: Decide whether we include the header here or not
}
