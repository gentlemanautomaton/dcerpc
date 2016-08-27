package ndr

import (
	"reflect"
	"sync"
)

// EncoderTypeCache represents a cache of types for which an RPC encoding
// engine has been compiled.
type EncoderTypeCache struct {
	mutex sync.RWMutex
	cache map[reflect.Type]EncOp
	// TODO: Allow compilers that are waiting on a particular type to pend its completion
	//pending map[reflect.Type][]
}

// NewEncoderTypeCache returns a new cache that is capable of storing types for
// which and RPC encoding engine has been compiled.
func NewEncoderTypeCache() *EncoderTypeCache {
	return &EncoderTypeCache{
		cache: make(map[reflect.Type]EncOp),
	}
}

// Add will add the given encoding operation to the cache for the given type.
func (c *EncoderTypeCache) Add(rt reflect.Type, op EncOp) {
	c.mutex.Lock()
	c.cache[rt] = op
	c.mutex.Unlock()
}

// Get returns the RPC encoding op for the requested type if it exists in the
// cache, or else nil.
func (c *EncoderTypeCache) Get(rt reflect.Type) (op EncOp) {
	c.mutex.RLock()
	op = c.cache[rt]
	c.mutex.RUnlock()
	return
}
