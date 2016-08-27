package ndr64

import (
	"io"
	"reflect"
	"sync"

	"github.com/gentlemanautomaton/dcerpc/formatlabel"
	"github.com/gentlemanautomaton/dcerpc/ndr"
)

// TODO: Store a separate type cache for each supported transfer syntax (NDR and NDR64)
var encTypeCache = ndr.NewEncoderTypeCache()

// Encoder encodes Go types as NDR64 data and transmits them via an underlying
// io.Writer..
type Encoder struct {
	mutex sync.Mutex
	w     ndr.Writer
}

// NewEncoder returns a new encoder that transmits nNDR64-encoded values on the
// given io.Writer. Writes are guarded by a mutex and are written atomicaly.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: ndr.NewWriter(w, formatlabel.LEAIEEE),
	}
}

// Encode encodes the given value in NDR64 and transmits the encoded value on
// the underlying io.Writer.
func (enc *Encoder) Encode(v interface{}) error {
	return enc.EncodeValue(reflect.ValueOf(v))
}

// EncodeValue encodes the given value in NDR64 and transmits the encoded value
// on the underlying io.Writer.
func (enc *Encoder) EncodeValue(v reflect.Value) error {
	//enc.cache.RLock()

	op := encTypeCache.Get(v.Type())
	// TODO: Add pending mechanism to avoid duplication of effort
	if op == nil {
		op = ndr.EncOpFor(v.Type()) // TODO: Replace with ndr64 function call
		encTypeCache.Add(v.Type(), op)
	}

	s := ndr.State{} // FIXME: Figure out how to handle state.

	enc.mutex.Lock()
	op(enc.w, &s, v) // TOOD: Return error from op?
	enc.mutex.Unlock()
	return nil
}
