package types

import "strings"

//type FieldAttrType int

const (
//FieldAttr
)

// FieldAttrList represents a list of field attributes.
type FieldAttrList []FieldAttr

// Contains returns true if the field attribute list contains one or more of
// the given field attribute types.
func (list FieldAttrList) Contains(anyOf ...string) bool {
	for i := 0; i < len(list); i++ {
		for t := 0; t < len(anyOf); t++ {
			if list[i].Type == anyOf[t] {
				return true
			}
		}
	}

	return false
}

// IsConformant returns true if the field attributes indicate a conformant field.
func (list FieldAttrList) IsConformant() bool {
	return list.Contains("min_is", "max_is", "size_is")
}

// IsVarying returns true if the field attributes indicate a varying field.
func (list FieldAttrList) IsVarying() bool {
	return list.Contains("first_is", "last_is", "length_is")
}

// FieldAttr represents a field attribute.
type FieldAttr struct {
	Type  string
	Value string
}

// ParseFieldAttrList parses the given field attribute list IDL string and
// returns the parsed data as a FieldAttrList.
func ParseFieldAttrList(attrs string) (output FieldAttrList) {
	values := strings.Split(attrs, ",")
	if len(values) == 0 {
		return
	}
	output = make([]FieldAttr, 0, len(values))
	for i := 0; i < len(values); i++ {
		attr, ok := ParseFieldAttr(strings.TrimSpace(values[i]))
		if ok {
			output = append(output, attr)
		}
		// FIXME: Handle errors
	}
	return
}

// ParseFieldAttr parses the given field attribute IDL string and
// returns the parsed data as a FieldAttr.
func ParseFieldAttr(attr string) (value FieldAttr, ok bool) {
	var t, v string
	switch attr {
	case "ignore":
		t = "ignore"
	default:
		t, v = parseParenthetical(attr)
	}
	value = FieldAttr{t, v}
	ok = (t != "")
	return
}

func parseParenthetical(p string) (typ, value string) {
	p1 := strings.Index(p, "(")
	if p1 < 0 {
		return
	}

	p2 := strings.LastIndex(p, ")")
	if p2 < 0 || p2 <= p1+1 {
		return
	}

	typ = p[0:p1]
	value = p[p1+1 : p2]

	return
}
