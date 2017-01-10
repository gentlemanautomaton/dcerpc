package ndr

// SliceSubset represents the offset and count of a single dimension of an
// N-dimensional varying array, which is stored as a slice.
type SliceSubset struct {
	Offset int
	Count  int
}
