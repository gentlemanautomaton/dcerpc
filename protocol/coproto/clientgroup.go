package coproto

import "sync"

type clientGroupKey struct {
	address string // FIXME: Figure out how we're storing primary server addresses
	id      uint32 // Association group ID
}

// ClientGroup manages the client side of an association group.
//
// Clients that share a common protocol tower and communicate with the same
// server are placed into the same client group.
type ClientGroup struct {
	id uint32

	mutex   sync.RWMutex
	clients []*Client
	active  uint // How many active contexts are there?
}

// add will add the client to the group.
func (group *ClientGroup) add(client *Client) {
	group.mutex.Lock()
	defer group.mutex.Unlock()
	group.clients = append(group.clients, client)
	return
}

// remove will remove client from the group.
func (group *ClientGroup) remove(client *Client) {
	group.mutex.Lock()
	defer group.mutex.Unlock()
	for i := 0; i < len(group.clients); i++ {
		if group.clients[i] == client {
			group.clients = append(group.clients[:i], group.clients[i+1:]...)
		}
	}
}

// activate will increment the number of active contexts within the client
// association group.
func (group *ClientGroup) activate() {
	group.mutex.Lock()
	defer group.mutex.Unlock()
	group.active++
}

// deactivate will decrement the number of associations within the client
// association group.
func (group *ClientGroup) deactivate() {
	group.mutex.Lock()
	defer group.mutex.Unlock()
	group.active--
}
