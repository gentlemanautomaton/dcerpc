package uuid

// TODO: Decide whether to continue rolling our own package for this here or
// instead to use one of the existing UUID packages floating around.

// UUID represents a 16 byte universally unique identifier as used by the
// Distributed Computing Environment.
type UUID [16]byte
