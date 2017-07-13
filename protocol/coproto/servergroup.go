package coproto

import "sync"

// ServerGroup manages the server side of an association group.
type ServerGroup struct {
	id uint

	mutex   sync.RWMutex
	servers []*Server
	active  uint // How many active contexts are there?
}
