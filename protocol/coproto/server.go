package coproto

import "sync"

// Server is a connection-oriented protocol server. It represents the server
// side of an RPC association.
//
// Each server is capabale of handling one RPC call at a time.
type Server struct {
	mutex sync.Mutex
	group *ServerGroup
}
