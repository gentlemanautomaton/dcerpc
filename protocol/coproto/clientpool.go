package coproto

import "sync"

// ClientPool manages a pool of clients.
type ClientPool struct {
	mutex  sync.RWMutex
	groups map[clientGroupKey]*ClientGroup
}

// Allocate will return an existing client from the pool if one is available,
// otherwise it will attempt to allocate a new client and return it.
//
// TODO: Receive and process the endpoint address as a parameter.
func (pool *ClientPool) Allocate() (client *Client) {
	// FIXME: Write this
	return nil
}
