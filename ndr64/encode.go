package ndr64

import (
	"reflect"

	"github.com/gentlemanautomaton/dcerpc/idl/types"
	"github.com/gentlemanautomaton/dcerpc/ndr"
)

type encInstr struct {
	op    ndr.EncOp
	index []int
}

// EncOpForSlice returns an NDR encoding function for the given type, which
// must be a slice.
func EncOpForSlice(rt reflect.Type) ndr.EncOp {
	// FIXME: Handle multiple dimensions
	elemOp := EncOpFor(rt.Elem())
	return func(w ndr.Writer, s *ndr.State, v reflect.Value) {
		w.WriteUint64(0)                 // Varying array offset, always zero in our case
		w.WriteUint64((uint64)(v.Len())) // Varying array length, in number of elements
		for e := 0; e < v.Len(); e++ {
			elemOp(w, s, v.Index(e))
		}
	}
}

// EncOpForConformantStruct returns an NDR64 conformant data encoding function
// for the given type, which must be a struct.
func EncOpForConformantStruct(rt reflect.Type) (ndr.EncOp, []int) {
	// TODO: Figure out a decent implementation for this in the ndr package
	//       and then create the appropriate implementation here.
	return ndr.EncNoop, nil
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
	if op := ndr.EncOpForPrimitive(rf.Type); op != nil {
		return op
	}

	attrs := types.ParseFieldAttrList(rf.Tag.Get("idl"))

	switch rf.Type.Kind() {
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
	if op := ndr.EncOpForPrimitive(rt); op != nil {
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
