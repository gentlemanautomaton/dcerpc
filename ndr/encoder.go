package ndr

import (
	"errors"
	"io"
)

// Encoder is capable of encoding Go types as Network Data Representation data.
type Encoder struct {
	w     Writer
	label FormatLabel
	// TODO: Add type map
}

// NewEncoder returns a new Encoder that writes to the given io.Writer with
// the encoding represented by the provided format label.
func NewEncoder(w io.Writer, label FormatLabel) (enc *Encoder, err error) {
	enc = &Encoder{
		w:     NewWriter(w, label),
		label: label,
	}
	if enc.w == nil {
		return nil, errors.New("Invalid format label")
	}
	return
}
