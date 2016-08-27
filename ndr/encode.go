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

// EncOp represents a compiled NDR encoding operation for a particular type.or
// field.
type EncOp func(w Writer, s *State, v reflect.Value)

// EncNoop is an NDR encoding function that does nothing..
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

}

// EncVaryingString is an NDR encoding function for varying strings.
func EncVaryingString(w Writer, s *State, v reflect.Value) {
	w.WriteUint32(0)                 // Varying string offset, always zero in our case
	w.WriteUint32((uint32)(v.Len())) // Varying string length, in number of characters
}

// EncOpForConformantStruct returns an NDR conformant data encoding function
// for the given type, which must be a struct.
func EncOpForConformantStruct(rt reflect.Type) (EncOp, []int) {
	last := rt.NumField() - 1
	if last >= 0 {
		f := rt.Field(last)
		if IsConformantField(f) {
			return EncOpForConformantField(f)
		}
	}
	return EncNoop, nil
}

// EncOpForConformantSlice returns an NDR conformant data encoding function for
// the given type, which must be a slice.
func EncOpForConformantSlice(rt reflect.Type, attrs types.FieldAttrList) EncOp {
	return nil
}

// EncOpForConformantField returns an NDR conformant data encoding function for
// the given field, which must be a conformant slice or a conformant struct.
func EncOpForConformantField(rf reflect.StructField) (EncOp, []int) {
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
	}
	if attrs.IsConformant() {
		// TODO: Look for conformant attribute data
		return nil, nil
	}
	// TODO: Evaluate struct members to find the conformant member
	return nil, nil
}

// EncOpForArray returns an NDR encoding functiong for the given type, which
// must be an array.
func EncOpForArray(rt reflect.Type) EncOp {
	// FIXME: Handle multiple dimensions
	count := rt.Len()
	elemOp := EncOpFor(rt.Elem())
	return func(w Writer, s *State, v reflect.Value) {
		for e := 0; e < count; e++ {
			elemOp(w, s, v.Index(e))
		}
	}
}

// EncOpForSlice returns an NDR encoding function for the given type, which
// must be a slice.
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

// EncOpForStruct returns an NDR encoding function for the given type, which
// must be a struct.
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

// EncOpForStructAlignment returns an NDR encoding function for aligning the
// given type, which must be a struct.
func EncOpForStructAlignment(rt reflect.Type) EncOp {
	return nil
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
