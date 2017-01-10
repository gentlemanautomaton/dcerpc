package ndr

import "fmt"

// Compile-time encoding error codes
const (
	MissingIDLFieldRef = 1000 + iota
)

// Run-time encoding error codes
const (
	FirstLessThanMin = 2000 + iota
	LastLessThanMin
	FirstGreaterThanLast
	NegativeSize
	NegativeLength
)

// EncodingError represents an error encountered during NDR encoding.
type EncodingError struct {
	Code         int
	TypeName     string
	FieldName    string
	RefFieldName string
	Value        int
	Limit        int
}

func (e EncodingError) Error() string {
	switch e.Code {
	case FirstGreaterThanLast:
		return fmt.Sprintf("NDR encoder error: Type \"%s\" contains a varying array field \"%s\" with an invalid first index \"%d\" that is greater than its last index \"%d\"", e.TypeNanme, e.FieldName, e.Index, e.Limit)
	case MissingIDLFieldRef:
		return fmt.Sprintf("NDR encoder error: Type \"%s\" does not contain the \"%s\" field, which was referenced by the IDL attributes of the \"%s\" field.", e.TypeName, e.RefFieldName, e.FieldName)
	default:
		return "Unknown NDR encoding error"
	}
}

// NewEncodingError returns an error for the given error code, type name, field
// name and referenced field name.
func NewEncodingError(code int, typeName, fieldName, refFieldName string, value, limit int) error {
	return &EncodingError{
		Code:         code,
		TypeName:     typeName,
		FieldName:    fieldName,
		RefFieldName: refFieldName,
		Value:        value,
		Limit:        limit,
	}
}
