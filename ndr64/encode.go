package ndr64

import (
	"reflect"

	"github.com/gentlemanautomaton/dcerpc/idl/types"
	"github.com/gentlemanautomaton/dcerpc/ndr"
)

// EncOpForSlice returns an NDR encoding function for the given type, which
// must be a slice.
func EncOpForSlice(rt reflect.Type) ndr.EncOp {
	// FIXME: Handle multiple dimensions
	elemOp := EncOpFor(rt.Elem())
	return func(w ndr.Writer, s *ndr.State, v reflect.Value) {
		w.WriteUint32(0)                 // Varying array offset, always zero in our case
		w.WriteUint32((uint32)(v.Len())) // Varying array length, in number of elements
		for e := 0; e < v.Len(); e++ {
			elemOp(w, s, v.Index(e))
		}
	}
}

// EncOpForStruct returns an NDR encoding function for the given type, which
// must be a struct.
func EncOpForStruct(rt reflect.Type) ndr.EncOp {
	engine := make([]encInstr, 0, rt.NumField())
	// TODO: Add alignment op?
	if ndr.IsConformantStruct(rt) {
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
	return func(w ndr.Writer, s *ndr.State, v reflect.Value) {
		for i := 0; i < len(engine); i++ {
			instr := &engine[i]
			field := v.FieldByIndex(instr.index)
			instr.op(w, s, field)
		}
	}
}

// EncOpForField returns an NDR encoding function for the given field.
func EncOpForField(rf reflect.StructField) ndr.EncOp {
	attrs := types.ParseFieldAttrList(rf.Tag.Get("idl"))

	switch rf.Type.Kind() {
	case reflect.Bool:
		return ndr.EncBool
	case reflect.Int8:
		return ndr.EncInt8
	case reflect.Uint8:
		return ndr.EncUint8
	case reflect.Int16:
		return ndr.EncInt16
	case reflect.Uint16:
		return ndr.EncUint16
	case reflect.Int32:
		return ndr.EncInt32
	case reflect.Uint32:
		return ndr.EncUint32
	case reflect.Int64:
		return ndr.EncInt64
	case reflect.Uint64:
		return ndr.EncUint64
	case reflect.Array:
		return ndr.EncOpForArray(rf.Type)
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

// EncOpFor returns an NDR encoding function for the given type.
func EncOpFor(rt reflect.Type) ndr.EncOp {
	// TODO: Figure out a good workaround for specifying attributes for non-fields
	//       Perhaps they could be namelessly composed into containing structs?
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
