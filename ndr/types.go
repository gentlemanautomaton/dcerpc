package ndr

import "reflect"

// SliceSubset represents the offset and count of a single dimension of an
// N-dimensional varying array, which is stored as a slice.
type SliceSubset struct {
	Offset int // Varying array offset, in number of elements
	Count  int // Varying array length, in number of elements
}

// SlicePosition is used internally when working with multi-dimensional slices.
type SlicePosition struct {
	index        int // current traversal index of this dimension, not including the offset
	parent       reflect.Value
	parentLength int // Actual length of parent slice
}
