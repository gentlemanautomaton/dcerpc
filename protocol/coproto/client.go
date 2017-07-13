package coproto

import "sync"

// Client is a connection-oriented protocol client. It represents the client
// side of an RPC association.
//
// Each client is capabale of making one RPC call at a time.
type Client struct {
	mutex sync.Mutex
	group *ClientGroup
}

// Close will release any resources allocated by the client.
//
// The client will remove itself from the group when it is closed.
func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.group == nil {
		return // Already closed
	}
	c.group.remove(c)
}

// Group returns the group that the client is a member of.
func (c *Client) Group() *ClientGroup {
	return c.group
}
