package copdu

// AlterContext represents an alter_context PDU in the connection-oriented
// protocol. It is sent from the client to the server. Its format is identical
// to Bind.
type AlterContext Bind
