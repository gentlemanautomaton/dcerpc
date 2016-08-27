package ndr

import (
	"errors"
	"io"
	"reflect"
	"sync"

	"github.com/gentlemanautomaton/dcerpc/formatlabel"
)

// TODO: Store a separate type cache for each supported transfer syntax (NDR and NDR64)
var encTypeCache = NewEncoderTypeCache()

// Encoder is capable of encoding Go types as Network Data Representation data.
type Encoder struct {
	mutex  sync.Mutex
	w      Writer
	format formatlabel.Format
	// TODO: Add type map
}

// NewEncoder returns a new Encoder that writes to the given io.Writer with
// the encoding represented by the provided format label.
func NewEncoder(w io.Writer, format formatlabel.Format) (enc *Encoder, err error) {
	enc = &Encoder{
		w:      NewWriter(w, format),
		format: format,
	}
	if enc.w == nil {
		return nil, errors.New("Invalid format label")
	}
	return
}

// Encode encodes the given value in NDR and transmits the encoded value on the
// underlying io.Writer.
func (enc *Encoder) Encode(v interface{}) error {
	return enc.EncodeValue(reflect.ValueOf(v))
}

// EncodeValue encodes the given value in network data representation format and
// writes it to the Encoder's underlying io.Writer.
func (enc *Encoder) EncodeValue(v reflect.Value) error {
	//enc.cache.RLock()

	op := encTypeCache.Get(v.Type())
	// TODO: Add pending mechanism to avoid duplication of effort
	if op == nil {
		op = EncOpFor(v.Type())
		encTypeCache.Add(v.Type(), op)
	}

	s := State{} // FIXME: Figure out how to handle state.

	enc.mutex.Lock()
	op(enc.w, &s, v) // TOOD: Return error from op?
	enc.mutex.Unlock()
	return nil
}

func (enc *Encoder) buildType() {

}
