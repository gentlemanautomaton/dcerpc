package ndr

import "sync"

// State represents the shared encoder and decoder state when communicating
// via NDR.
type State struct {
	mutex sync.RWMutex // TODO: Use different mutexes for different properties
	// ptrToRef maps pointers to referent identifiers. Pointers are cast to
	// interface{} to avoid the use of unsafe.Pointer.
	ptrToRef map[interface{}]uint64
	// refToPtr maps referent identifiers to pointers. Pointers are cast to
	// interface{} to avoid the use of unsafe.Pointer.
	refToPtr map[uint64]interface{}
	// id is the last referent ID to be allocated
	id uint64
	// errors is the set of errors encountered during encoding or decoding
	errors []error
}

// NewState initializes a new encoder/decoder state and returns it.
func NewState() *State {
	return &State{
		ptrToRef: make(map[interface{}]uint64),
		refToPtr: make(map[uint64]interface{}),
	}
}

// Register assigns a referent identifier to the given value and returns it. If
// the value was registered previously, the previously assigned referent
// identifier is returned.
//
// The provided value must be a pointer type.
//
// The returned referent identifier is a monotonically increasing 64-bit value.
func (s *State) Register(v interface{}) (refID uint64) {
	// First pass with read lock
	s.mutex.RLock()
	refID, ok := s.ptrToRef[v]
	s.mutex.RUnlock()
	if ok {
		return
	}

	// Second pass with write lock
	s.mutex.Lock()
	defer s.mutex.Unlock()
	refID, ok = s.ptrToRef[v]
	if ok {
		return
	}
	s.id++
	s.ptrToRef[v] = s.id
	s.refToPtr[s.id] = v
	refID = s.id
	return
}

// AddError adds the given error to the list of errors encountered in the
// current encoding or decoding session.
func (s *State) AddError(err error) {
	s.mutex.Lock()
	s.errors = append(s.errors, err)
	s.mutex.Unlock()
}
