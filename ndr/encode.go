package ndr

import (
	"reflect"

	"github.com/gentlemanautomaton/dcerpc/idl/types"
)

type encEngine struct {
	alignment int
	instr     []encInstr
}

type encInstr struct {
	op    EncOp
	index []int
}

type encField struct {
	Name  string
	Index []int
}

// EncOp represents a compiled NDR encoding operation for a particular type or
// field.
type EncOp func(w Writer, s *State, v reflect.Value)

// EncNoop is an NDR encoding function that does nothing.
func EncNoop(w Writer, s *State, v reflect.Value) {}

// EncBytes is an NDR encoding function for a byte slice.
func EncBytes(w Writer, s *State, v reflect.Value) {
	w.Write(v.Bytes())
}

// EncBool is an NDR encoding function for a bool.
func EncBool(w Writer, s *State, v reflect.Value) {
	w.WriteBool(v.Bool())
}

// EncInt8 is an NDR encoding function for an int8.
func EncInt8(w Writer, s *State, v reflect.Value) {
	w.WriteInt8((int8)(v.Int()))
}

// EncUint8 is an NDR encoding function for a uint8.
func EncUint8(w Writer, s *State, v reflect.Value) {
	w.WriteUint8((uint8)(v.Uint()))
}

// EncInt16 is an NDR encoding function for an int16.
func EncInt16(w Writer, s *State, v reflect.Value) {
	w.WriteInt16((int16)(v.Int()))
}

// EncUint16 is an NDR encoding function for a uint16.
func EncUint16(w Writer, s *State, v reflect.Value) {
	w.WriteUint16((uint16)(v.Uint()))
}

// EncInt32 is an NDR encoding function for an int32.
func EncInt32(w Writer, s *State, v reflect.Value) {
	w.WriteInt32((int32)(v.Int()))
}

// EncUint32 is an NDR encoding function for a uint32.
func EncUint32(w Writer, s *State, v reflect.Value) {
	w.WriteUint32((uint32)(v.Uint()))
}

// EncInt64 is an NDR encoding function for an int64.
func EncInt64(w Writer, s *State, v reflect.Value) {
	w.WriteInt64(v.Int())
}

// EncUint64 is an NDR encoding function for a uint64.
func EncUint64(w Writer, s *State, v reflect.Value) {
	w.WriteUint64((uint64)(v.Uint()))
}

// EncString is an NDR encoding function for a string.
func EncString(w Writer, s *State, v reflect.Value) {
	w.WriteString(v.String())
}

// EncVaryingString is an NDR encoding function for varying strings.
func EncVaryingString(w Writer, s *State, v reflect.Value) {
	w.WriteUint32(0)                 // Varying string offset, always zero in our case
	w.WriteUint32((uint32)(v.Len())) // Varying string length, in number of characters
	w.WriteString(v.String())
}

// EncOpForArray returns an NDR encoding function for the given type, which
// must be an array.
func EncOpForArray(rt reflect.Type) EncOp {
	e1 := rt.Elem()
	if e1.Kind() != reflect.Array {
		return EncOpForArray1D(rt.Len(), EncOpFor(e1))
	}
	e2 := e1.Elem()
	if e2.Kind() != reflect.Array {
		return EncOpForArray2D(rt.Len(), e1.Len(), EncOpFor(e2))
	}
	e3 := e2.Elem()
	if e3.Kind() != reflect.Array {
		return EncOpForArray3D(rt.Len(), e1.Len(), e2.Len(), EncOpFor(e3))
	}
	e4 := e3.Elem()
	return EncOpForArray4D(rt.Len(), e1.Len(), e2.Len(), e3.Len(), EncOpFor(e4))
}

// EncOpForArray1D returns an NDR encoding function for a one-dimensional array
// with the given length and element encoding function.
func EncOpForArray1D(length int, elemOp EncOp) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		for i := 0; i < length; i++ {
			elemOp(w, s, v.Index(i))
		}
	}
}

// EncOpForArray2D returns an NDR encoding function for a two-dimensional array
// with the given lengths and element encoding function.
func EncOpForArray2D(len1, len2 int, elemOp EncOp) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				elemOp(w, s, v.Index(i).Index(j))
			}
		}
	}
}

// EncOpForArray3D returns an NDR encoding function for a three-dimensional
// array with the given lengths and element encoding function.
func EncOpForArray3D(len1, len2, len3 int, elemOp EncOp) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for k := 0; k < len3; k++ {
					elemOp(w, s, v.Index(i).Index(j).Index(k))
				}
			}
		}
	}
}

// EncOpForArray4D returns an NDR encoding function for a three-dimensional
// array with the given lengths and element encoding function.
func EncOpForArray4D(len1, len2, len3, len4 int, elemOp EncOp) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for k := 0; k < len3; k++ {
					for m := 0; m < len4; m++ {
						elemOp(w, s, v.Index(i).Index(j).Index(k).Index(m))
					}
				}
			}
		}
	}
}

// EncOpForSlice returns an NDR encoding function for the given type, which
// must be a slice. The encoding function will encode the slice will as a
// varying array.
func EncOpForSlice(rt reflect.Type) EncOp {
	// FIXME: Handle multiple dimensions
	elemOp := EncOpFor(rt.Elem())
	return func(w Writer, s *State, v reflect.Value) {
		w.WriteUint32(0)                 // Varying array offset, always zero in our case
		w.WriteUint32((uint32)(v.Len())) // Varying array length, in number of elements
		for e := 0; e < v.Len(); e++ {
			elemOp(w, s, v.Index(e))
		}
	}
}

// EncSliceElements encodes all of the slice elements with the given slice
// element encoding function. The slice may be multi-dimensional. Only the
// subsets specified for each dimension will be encoded.
//
// If any subset exceeds the boundary of the actual slice data, zero values
// appropriate for the element type will be generated to fill in the place of
// the missing range. This is done to avoid encoding malformed data.
func EncSliceElements(w Writer, s *State, v reflect.Value, subsets []SliceSubset, elemOp EncOp) {
	type position struct {
		index        int // current traversal index of this dimension, not including the offset
		parent       reflect.Value
		parentLength int // Actual length of parent slice
	}
	dimensions := len(subsets)
	maxDepth := dimensions - 1
	depth := 0
	stack := make([]position, 0, dimensions)
	stack = append(stack, position{0, v, v.Len()})
	pos, subset := &stack[depth], &subsets[depth]
	for {
		// This loop stops at each slice in the N-dimensional set of slices (which
		// is effectively a tree of slices). It starts at the root slice.

		// Traverse down to the next leaf slice
		for depth < maxDepth {
			depth++
			if pos.parent.IsValid() && pos.index < pos.parentLength {
				p := pos.parent.Index(pos.index)
				stack = append(stack, position{0, p, p.Len()})
			} else {
				stack = append(stack, position{0, reflect.Value{}, 0})
			}
			pos, subset = &stack[depth], &subsets[depth]
		}

		// Write all of the elements in the leaf slice
		start, end := subset.Offset, subset.Offset+subset.Count
		if pos.parent.IsValid() {
			end1, end2 := pos.parentLength, end
			if end < end2 {
				//
			}
			// Write values that are present
			if pos.parentLength >= subset.Offset+subset.Count {

			}
			// Write values that are absent
		} else {
			// Write zero elements entirely
			for i := start; i < end; i++ {
				// FIXME: Add error to state
				// TODO: Accept a zero-value element function and call it here?
				elemOp(w, s, reflect.Value{})
			}
		}
		for i, end := subset.Offset, subset.Offset+subset.Count; i < end; i++ {
			// FIXME: Perform boundary checks on the actual data
			if i < pos.parentLength {
				// FIXME: Add error to state
				// TODO: Accept a zero-value element function and call it here?
				elemOp(w, s, reflect.Value{})
			} else {
				elemOp(w, s, pos.parent.Index(i))
			}
		}

		// Traverse upward and forward
		for {
			depth--
			if depth < 0 {
				return
			}
			stack = stack[0 : len(stack)-1]
			pos, subset = &stack[depth], &subsets[depth]
			pos.index++
			if pos.index < subset.Count {
				break
			}
		}
	}
}

// EncOpForStruct returns an NDR encoding function for the given type, which
// must be a struct. If the struct contains conformant data it will be
// encoded appropriately.
//
// See section 14.3.6 of the DCE RPC publication for an overview of the
// struct encoding rules under NDR transfer syntax.
func EncOpForStruct(rt reflect.Type) EncOp {
	engine := make([]encInstr, 0, rt.NumField())
	// TODO: Add alignment op?
	if IsConformantStruct(rt) {
		op, index := EncOpForConformantStruct(rt)
		engine = append(engine, encInstr{
			op:    op,
			index: index,
		})
	}
	if alignmentOp := EncOpForStructAlignment(rt); alignmentOp != nil {
		engine = append(engine, encInstr{
			op: alignmentOp,
		})
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if op := EncOpForField(f); op != nil {
			engine = append(engine, encInstr{
				op:    op,
				index: f.Index,
			})
		}
	}
	return func(w Writer, s *State, v reflect.Value) {
		for i := 0; i < len(engine); i++ {
			instr := &engine[i]
			field := v.FieldByIndex(instr.index)
			instr.op(w, s, field)
		}
	}
}

func encOpForInstructions(engine []encInstr, alignment int) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		w.Align(alignment)
		for i := 0; i < len(engine); i++ {
			instr := &engine[i]
			field := v.FieldByIndex(instr.index)
			instr.op(w, s, field)
		}
	}
}

func encOpForConformantInstructions(engine []encInstr, alignment int, prealignmentCount int) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		for i := 0; i < prealignmentCount; i++ {
			instr := &engine[i]
			field := v.FieldByIndex(instr.index)
			instr.op(w, s, field)
		}
		w.Align(alignment)
		for i := prealignmentCount; i < len(engine); i++ {
			instr := &engine[i]
			field := v.FieldByIndex(instr.index)
			instr.op(w, s, field)
		}
	}
}

// EncOpForStructAlignment returns an NDR encoding function for aligning the
// given type, which must be a struct.
func EncOpForStructAlignment(rt reflect.Type) EncOp {
	return nil
}

// EncOpForPrimitive returns an NDR encoding function for the given type, if it
// represents an NDR primitive, otherwise it returns nil.
func EncOpForPrimitive(rt reflect.Type) EncOp {
	switch rt.Kind() {
	case reflect.Bool:
		return EncBool
	case reflect.Int8:
		return EncInt8
	case reflect.Uint8:
		return EncUint8
	case reflect.Int16:
		return EncInt16
	case reflect.Uint16:
		return EncUint16
	case reflect.Int32:
		return EncInt32
	case reflect.Uint32:
		return EncUint32
	case reflect.Int64:
		return EncInt64
	case reflect.Uint64:
		return EncUint64
	}
	return nil
}

// EncOpForStructConformance returns an NDR conformant data encoding function
// for the given type, which must be a struct.
func EncOpForStructConformance(rt reflect.Type) (EncOp, []int) {
	last := rt.NumField() - 1
	if last >= 0 {
		f := rt.Field(last)
		if IsConformantField(f) {
			return EncOpForFieldConformance(f)
		}
	}
	return EncNoop, nil
}

// EncOpForSliceConformance returns an encoding function for NDR conformance
// data that encodes the conformance offset for the given base type, and slice
// type which must be a
// slice.
func EncOpForSliceConformance(base reflect.Type, slice reflect.Type, attrs types.FieldAttrList) EncOp {
	min, hasMin := attrs.Lookup("min_is")
	max, hasMax := attrs.Lookup("max_is")
	size, hasSize := attrs.Lookup("size")
	switch {
	case hasMin && hasMax:
		minField, minOk := base.FieldByName(min)
		maxField, maxOk := base.FieldByName(max)
		if minOk && maxOk {
			return EncOpForMinMaxConformance(minField.Index, maxField.Index)
		}
		// FIXME: panic?
	case hasSize:
		sizeField, sizeOk := base.FieldByName(size)
		if sizeOk {
			return EncOpForSizeConformance(sizeField.Index)
		}
	case hasMin:
		minField, minOk := base.FieldByName(min)
		if minOk {
			return EncOpForMinConformance(minField.Index)
		}
	case hasMax:
		maxField, maxOk := base.FieldByName(max)
		if maxOk {
			return EncOpForMaxConformance(maxField.Index)
		}
	}
	// FIXME: panic?
	return EncNoop
}

func EncOpForMinMaxConformance(minFieldIndex []int, maxFieldIndex []int) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		min := v.FieldByIndex(minFieldIndex).Uint()
		max := v.FieldByIndex(maxFieldIndex).Uint()
		// FIXME: panic if min or max overflows a 32 bit integer?
		w.WriteUint32((uint32)(min))
		w.WriteUint32((uint32)(max))
	}
}

func EncOpForMinConformance(minFieldIndex []int) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		min := v.FieldByIndex(minFieldIndex).Uint()
		// FIXME: panic if min overflows a 32 bit integer?
		w.WriteUint32((uint32)(min))
	}
}

func EncOpForMaxConformance(maxFieldIndex []int) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		max := v.FieldByIndex(maxFieldIndex).Uint()
		// FIXME: panic if max overflows a 32 bit integer?
		w.WriteUint32((uint32)(max))
	}
}

func EncOpForSizeConformance(sizeFieldIndex []int) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		size := v.FieldByIndex(sizeFieldIndex).Uint()
		// FIXME: panic if size overflows a 32 bit integer?
		w.WriteUint32((uint32)(size))
	}
}

// EncOpForConformantField returns an NDR conformant data encoding function for
// the given field, which must be a conformant slice or a conformant struct.
func EncOpForConformantField(base reflect.Type, rf reflect.StructField) (EncOp, []int) {
	attrs := types.ParseFieldAttrList(rf.Tag.Get("idl"))
	switch rf.Type.Kind() {
	case reflect.Slice:
		if attrs.IsConformant() {
			// TODO: Evaluate min_is, max_is and size_is
			//op := encOpForConformantSlice(rf.Type, attrs)
			return func(w Writer, s *State, v reflect.Value) {
				return
			}, rf.Index
		}
	case reflect.Struct:
		if IsConformantStruct(rf.Type) {
			return
		}
	}
	if attrs.IsConformant() {
		// TODO: Look for conformant attribute data
		return nil, nil
	}
	// TODO: Evaluate struct members to find the conformant member
	return nil, nil
}

// EncOpForField returns an NDR encoding function for the given field.
func EncOpForField(rf reflect.StructField) EncOp {
	if op := EncOpForPrimitive(rf.Type); op != nil {
		return op
	}

	attrs := types.ParseFieldAttrList(rf.Tag.Get("idl"))

	switch rf.Type.Kind() {
	case reflect.Array:
		return EncOpForArray(rf.Type)
	case reflect.Slice:
		return EncOpForSlice(rf.Type)
	case reflect.String:
		if !attrs.IsConformant() && !attrs.IsVarying() {
			// Do something
		}
	case reflect.Struct:
		return EncOpForStruct(rf.Type)
	case reflect.Ptr:
		if attrs.Contains("ignore") {
			return nil
		}
	}
	return nil
}

// EncOpForSliceField returns an NDR encoding function for the given type, which
// must be a slice.
func EncOpForSliceField(base reflect.Type, slice reflect.Type, attrs types.FieldAttrList) EncOp {
	// FIXME: Handle multiple dimensions
	elemOp := EncOpFor(slice.Elem()) // FIXME: Handle embedded slices as multi-dimensional data?
	var (
		firstFieldName, hasFirst   = attrs.Lookup("first_is")
		lastFieldName, hasLast     = attrs.Lookup("last_is")
		lengthFieldName, hasLength = attrs.Lookup("length_is")
	)

	// FIXME: Support expressions?

	switch {
	case hasFirst && hasLast:
		firstField, firstOk := base.FieldByName(firstFieldName)
		lastField, lastOk := base.FieldByName(lastFieldName)
		if firstOk && lastOk {
			encOpForSliceWithFirstLast(firstField.Index, lastField.Index, elemOp)
		} else {
			if !firstOk {
				panic(NewEncodingError(MissingIDLFieldRef, base.Name(), slice.Name(), firstFieldName))
			}
			if !lastOk {
				panic(NewEncodingError(MissingIDLFieldRef, base.Name(), slice.Name(), lastFieldName))
			}
		}
	case hasFirst && hasLength:
	case hasLast:
	case hasLength:
	case hasLast && hasLength:
		panic("attribute list includes both last_is and length_is attributes, which are mutually exclusive")
	default:

		return func(w Writer, s *State, v reflect.Value) {
			w.WriteUint32(0)                 // Varying array offset, always zero in our case
			w.WriteUint32((uint32)(v.Len())) // Varying array length, in number of elements
			for e := 0; e < v.Len(); e++ {
				elemOp(w, s, v.Index(e))
			}
		}
	}
}

func encOpForSliceWithFirstLast(typeName string, sliceField, firstField, lastField encField, elemOp EncOp) EncOp {
	return func(w Writer, s *State, v reflect.Value) {
		first := v.FieldByIndex(firstField.Index).Uint()
		last := v.FieldByIndex(lastField.Index).Uint()
		// FIXME: panic if min or max overflows a 32 bit integer?
		// FIXME: Perform bounds checking?
		if first > last {
			s.AddError(NewEncodingError(FirstGreaterThanLast, typeName, sliceField.Name, firstField.Name, int(first), int(last)))
			w.WriteUint32(0)
			w.WriteUint32(0)
		} else {
			w.WriteUint32((uint32)(first))
			w.WriteUint32((uint32)(last - first + 1))
			for e := first; e < last; e++ {
				elemOp(w, s, v.Index(int(e)))
			}
		}
	}
}

// EncOpFor returns an NDR encoding function for the given type.
func EncOpFor(rt reflect.Type) EncOp {
	// TODO: Figure out a good workaround for specifying attributes for non-fields
	//       Perhaps they could be namelessly composed into containing structs?
	//       Alternatively: include empty struct types in into the struct that
	//                      signify behaviors.
	if op := EncOpForPrimitive(rt); op != nil {
		return op
	}

	switch rt.Kind() {
	case reflect.Array:
	case reflect.String:
		//if !attrs.IsConformant() && !attrs.IsVarying() {
		// Do something
		//}
	case reflect.Struct:
		return EncOpForStruct(rt)
	}
	return nil
}

// IsConformantStruct returns true if the given type is a conformant struct.
func IsConformantStruct(rt reflect.Type) bool {
	if rt.Kind() != reflect.Struct {
		return false
	}
	last := rt.NumField() - 1
	if last >= 0 {
		f := rt.Field(last)
		return IsConformantField(f)
	}
	return false
}

// IsConformantField returns true if the given field is conformant.
func IsConformantField(rf reflect.StructField) bool {
	switch rf.Type.Kind() {
	case reflect.Slice:
		attrs := types.ParseFieldAttrList(rf.Tag.Get("idl"))
		return attrs.IsConformant()
	case reflect.Struct:
		return IsConformantStruct(rf.Type)
	}
	return false
}

/*
func (enc *Encoder) encodeStruct(engine *encEngine, v reflect.Value) {
	for i := 0; i < len(engine.instr); i++ {
		instr := &engine.instr[i]
		field := v.FieldByIndex(instr.index)
		instr.op(instr, enc.w, field)
	}
}

func (enc *Encoder) encodeArray(engine *encEngine, v reflect.Value) {
	for i := 0; i < len(engine.instr); i++ {
		instr := &engine.instr[i]
		field := v.FieldByIndex(instr.index)
		instr.op(instr, enc.w, field)
	}
}
*/
