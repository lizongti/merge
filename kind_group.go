package merge

import "reflect"

type kindGroup int

const (
	kindGroupInvalid kindGroup = iota
	kindGroupContainer
	kindGroupRefer
	kindGroupValue
)

func getKindGroup(v reflect.Value) kindGroup {
	switch v.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array,
		reflect.Chan:
		return kindGroupContainer
	case reflect.Func, reflect.Interface, reflect.Pointer,
		reflect.UnsafePointer:
		return kindGroupRefer
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return kindGroupValue
	default:
		return kindGroupInvalid
	}
}
