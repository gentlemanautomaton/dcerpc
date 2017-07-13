package coproto

import "time"

const (
	// MaxBackoff is the maximum allowable time between retries.
	//
	// CONST_MAX_BACKOFF
	MaxBackoff = time.Minute

	// MaxResourceWait is the maximum allowable time for an association to be
	// allocated.
	//
	// CONST_MAX_RESOURCE_WAIT
	MaxResourceWait = time.Minute * 5

	// MinSupportedFragmentSize is the minimum size of PDU fragments that a
	// connection-oriented client or server must support.
	//
	// CONST_MUST_RCV_FRAG_SIZE
	MinSupportedFragmentSize = 1432
)
